"use strict";

const request = require("request");

const express = require("express");
const morgan = require("morgan");

// create a new express application
const app = express();

// add JSON request body parsing middleware
app.use(express.json());

// add request logging middleware
app.use(morgan("dev"));

const addr = process.env.ADDR || ":80";
const [host, port] = addr.split(":");
const portNum = parseInt(port);
if (isNaN(portNum)) {
	throw new Error("Port number is not a number.");
}

const mongo = process.env.MONGO || ":27017";
let url = 'mongodb://' + mongo + '/db';

const mongoose = require("mongoose");
const autoIncrement = require("mongoose-sequence")(mongoose);
const mongooseUniqueValidator = require("mongoose-unique-validator");

mongoose.connect(url, { useCreateIndex: true, useNewUrlParser: true });

const Schema = mongoose.Schema;
let petSchema = new Schema({
    _id: Number,
    name: {
        type: String,
        required: true,
        unique: true
    },
    gender: {
        type: String,
        required: true
    },
	breed: {
        type: String,
        required: true
    },
    age: Number,
    bio: String,
    petOwner: {
		id: Number,
		username: String,
		firstname: String,
		lastname: String,
		photourl: String
    },
    likeList: [Number]
});
petSchema.plugin(autoIncrement);
petSchema.plugin(mongooseUniqueValidator);
petSchema.set("toJSON", {
	transform: function (_, ret) {
		ret.id = ret._id;
		delete ret._id;
		delete ret.__v;
	}
});
let PetData = mongoose.model("Pet", petSchema);

// import rabbit 
const amqp = require('amqplib/callback_api');
const rabbit = process.env.RABBITADDR;
const queueName = "PetQueue";
let channel;
amqp.connect('amqp://' + rabbit + '/', function (err, conn) {
	if (err) {
		console.log("Failed to connect to Rabbit Instance from API Server.");
		process.exit(1);
	}

	conn.createChannel((err, ch) => {
		if (err) {
			console.log("Failed to create channel from API Server");
			process.exit(1);
		}

		ch.assertQueue(queueName, { durable: true });
		channel = ch;
		console.log(`Rabbit is running`);
	});
});

// gets information of the pet that is owned by the current user
app.get("/v1/pet", (req, res) => {
	let userJSON = req.get("X-User");
    if (userJSON) {
        let user = JSON.parse(userJSON);
        PetData.find({}).exec((err, pets) => {
            if (err) {
				return res.status(400).send("Error getting pets");
            } else {
                let pet = findOwnPet(pets, user);
                if (pet == undefined) {
                    return res.status(200).send(JSON.stringify(""));
                } else {
                    res.setHeader("Content-Type", "application/json");
                    return res.status(200).send(JSON.stringify(pet));
                }                
            }
        })
    } else {
		return res.status(401).send("User not authenticated");
	}
})

// creates a new pet profile for the current user
// responds with a copy of the newly-created pet
app.post("/v1/pet", (req, res) => {
	let userJSON = req.get("X-User");
    if (userJSON) {
        let user = JSON.parse(userJSON);
        PetData.find({}).exec((err, pets) => {
            if (err) {
                return res.status(400).send("Error getting pets");
            } else if (findOwnPet(pets, user) != undefined) {
                return res.status(400).send("Each user can only have one pet");
            } else if (req.body.name == undefined || req.body.name == "") {
                return res.status(400).send("Pet name is required");
            } else if (req.body.breed == undefined || req.body.breed == "") {
                return res.status(400).send("Pet breed is required");
            } else if (req.body.gender == undefined || req.body.gender == "") {
                return res.status(400).send("Pet gender is required");
            } else {
                let newPet = new PetData({
                    name: req.body.name,
                    gender: req.body.gender,
                    breed: req.body.breed,
                    age: req.body.age,
                    bio: req.body.bio,
                    petOwner: user,
                    likeList: []
                });
                newPet.save(err => {
                    if (err) {
                        return res.status(400).send("Error saving new pet");
                    } else {
                        const newMsg = {
                            "type": "pet-new",
                            "pet": newPet,
                            "userID": [user.id]
                        };
                
                        channel.sendToQueue(
                            queueName, 
                            new Buffer(JSON.stringify(newMsg)),
                            {persistent: true}
                        );

                        
                        res.setHeader("Content-Type", "application/json");
                        return res.status(201).send(JSON.stringify(newPet));
                    }
                });
            }
        });
	} else {
		return res.status(401).send("User not authenticated");
	}
})

// updates information of the pet owned by the current user
// responds with a copy of the newly-updated pet
app.patch("/v1/pet", (req, res) => {
	let userJSON = req.get("X-User");
    if (userJSON) {
        let user = JSON.parse(userJSON);
        PetData.find({}).exec((err, pets) => {
            if (err) {
				return res.status(400).send("Error getting pets");
            }
            let pet = findOwnPet(pets, user);
            if (pet == undefined) {
				return res.status(400).send("User has no pet");
            }
            PetData.update({ "_id": pet._id }, { $set: { name: req.body.name || pet.name, age: req.body.age || pet.age, bio: req.body.bio || pet.bio } }, err => {
                if (err) {
                    return res.status(500).send("Error updating");
                } else {
                    PetData.findOne({ "_id": pet._id }).exec((err, updated) => {
                        if (err) {
                            return res.status(500).send("Error updating");
                        } else {
                            res.setHeader("Content-Type", "application/json");
                            console.log(updated)
                            return res.status(200).send(JSON.stringify(updated));
                        }
                    });
                }
            })
        })
    } else {
		return res.status(401).send("User not authenticated");
    }
})

// deletes current user's pet profile
app.delete("/v1/pet", (req, res) => {
	let userJSON = req.get("X-User");
    if (userJSON) {
        let user = JSON.parse(userJSON);
        PetData.find({}).exec((err, pets) => {
            if (err) {
				return res.status(400).send("Error getting pets");
            }
            let pet = findOwnPet(pets, user);
            if (pet == undefined) {
				return res.status(400).send("User has no pet");
            }
            PetData.remove({ "_id": pet._id }, err => {
                if (err) {
                    return res.status(500).send("Error removing Pet")
                } else {
                    PetData.find({}).exec((err, pets) => {
                        if (err) {
                            return res.status(400).send("Error getting pets");
                        }

                        // remove the deleted pet from other pets' like lists
                        for (let p of pets) {
                            if (p.likeList.includes(pet._id)) {
                                let list = p.likeList
                                list.splice(list.indexOf(pet._id), 1)
                                PetData.update({ "_id": p._id }, { $set: { likeList: list } }, err => {
                                    if (err) {
                                        return res.status(500).send("Error deleting pet from other pet's liked list");
                                    } 
                                })
                            }
                        }
                        res.setHeader("Content-Type", "text/plain");
                        return res.status(200).send("Delete successful");
                    })
                }
            })
        }) 
    } else {
		return res.status(401).send("User not authenticated");
    }
})

// gets a list of pets for matching
app.get("/v1/pet/matching", (req, res) => {
	let userJSON = req.get("X-User");
    if (userJSON) {
        let user = JSON.parse(userJSON);
        PetData.find({}).exec((err, pets) => {
            if (err) {
				return res.status(400).send("Error getting pets");
			}
            
            // matching list should not contain own pet
            // should only contain pets of the same breed
            let pet = findOwnPet(pets, user);
            let list = [];
            for (let p of pets) {                
                if (p != pet && p.breed == pet.breed) {
                    list.push(p);
                }
            }

			res.setHeader("Content-Type", "application/json");
			return res.status(200).send(JSON.stringify(list));
		});
    } else {
        return res.status(401).send("User not authenticated");
    }
});

// likes a pet with the associated petID
// stores the petID in the like list of the pet that is owned by the current user
app.post("/v1/pet/:petID", (req, res) => {
	let userJSON = req.get("X-User");
    if (userJSON) {
        let user = JSON.parse(userJSON);
        PetData.find({}).exec((err, pets) => {
            if (err) {
				return res.status(400).send("Error getting pets");
			}
            
            let pet = findOwnPet(pets, user);
            // prevent duplicates in like lists
            for (let i = 0; i < pet.likeList.length; i++) { 
                if (pet.likeList[i] == req.params.petID) {
                    pet.likeList.splice(i, 1);
                    break;
                }
            }
            pet.likeList.push(req.params.petID);
            PetData.update({ "_id": pet._id }, { $set: {likeList: pet.likeList} }, err => {
                if (err) {
                    res.setHeader("Content-Type", "text/plain");
                    return res.status(500).send("Error liking pet");
                }

                // check if there is a match
                PetData.findOne({"_id": req.params.petID}).exec((err, p) => {
                    if (err || p == undefined) {
                        return res.status(400).send("Error finding liked pet");
                    }
                    for (let like of p.likeList) {
                        if (like == pet._id){
                            let chName = Date.now() + Math.random();
                            request.post({
                                headers: {
                                    "Content-Type": "application/json",
                                },
                                url: "https://api.demitu.me/v1/channels",
                                body: JSON.stringify({
                                    name: chName.valueOf(),  // generate a unique name each time
                                    private: true,
                                    members: [{"id": user.id}, {"id": p.petOwner.id}]
                                }),
                            }, (err, _, data) => {
                                if (err) {
                                    return res.status(400).send("Error creating channel");
                                } else {
                                    const newMsg = {
                                        "type": "Successfully matched",
                                        "channel": JSON.parse(data),
                                        "userIDs": [user.id, p.petOwner.id]
                                    };
                                    channel.sendToQueue(
                                        queueName, 
                                        new Buffer(JSON.stringify(newMsg)),
                                        {persistent: true}
                                    );
                                }
                            })
                        }
                    }
                    res.setHeader("Content-Type", "application/json");
                    return res.status(201).send(JSON.stringify(pet));
                })
            });
        });
    } else {
        return res.status(401).send("User not authenticated");
    }
});

// deletes the petID from the like list of the pet that is owned by the current user
app.delete("/v1/pet/:petID", (req, res) => {
	let userJSON = req.get("X-User");
    if (userJSON) {
        let user = JSON.parse(userJSON);
        PetData.find({}).exec((err, pets) => {
            if (err) {
				return res.status(400).send("Error getting pets");
			} else {
                let pet = findOwnPet(pets, user);
                // go through the like list of the pet that is owned by the current user            
                for (let i = 0; i < pet.likeList.length; i++) { 
                    if (pet.likeList[i] == req.params.petID) {
                        // delete the pet with the associated petID
                        pet.likeList.splice(i, 1);  // CHECK
                    }
                }
                PetData.update({ "_id": pet._id }, { $set: { likeList: pet.likeList } }, err => {
                    if (err) {
                        res.setHeader("Content-Type", "text/plain")
                        return res.status(500).send("Error unliking pet");
                    } else {
                        res.setHeader("Content-Type", "text/plain");
                        return res.status(200).send("Delete successful");
                    }
                });
            }
        });
    } else {
        return res.status(401).send("User not authenticated");
    }
});

// start the server listening on host:port
app.listen(portNum, host, () => {
	// callback is executed once server is listening
	console.log(`server is listening at http://${addr}...`);
});

// find the pet that is owned by the current user
function findOwnPet(pets, user) {
    for (let p of pets) {                
        if (p.petOwner.id == user.id) {
            return p;
        }
    }
}
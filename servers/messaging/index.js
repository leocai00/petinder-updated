"use strict";

const express = require("express");
const morgan = require("morgan");

// import rabbit 
let amqp = require('amqplib/callback_api');
const queueName = "MsgQueue";

// create a new express application
const app = express();

// add JSON request body parsing middleware
app.use(express.json());

// add the request logging middleware
app.use(morgan("dev"));

const addr = process.env.ADDR || ":80";

const rabbit = process.env.RABBITADDR;

const [host, port] = addr.split(":");
const portNum = parseInt(port);
if (isNaN(portNum)) {
	throw new Error("Port number is not a number.");
}

let channel;
amqp.connect('amqp://' + rabbit + '/', function (err, conn) {
	if (err) {
		console.log("Failed to connect to Rabbit Instance from API Server.");
		process.exit(1);
	}

	conn.createChannel(function (err, ch) {

		if (err) {
			console.log("Failed to create channel from API Server");
			process.exit(1);
		}

		ch.assertQueue("MsgQueue", { durable: true });
		channel = ch;

		console.log(`Rabbit is running`);
	});
});

const mongo = process.env.MONGO || ":27017";
let url = 'mongodb://' + mongo + '/db';

const mongoose = require("mongoose");
const autoIncrement = require("mongoose-sequence")(mongoose);
const ObjectId = require('mongodb').ObjectId;

mongoose.connect(url, { useCreateIndex: true, useNewUrlParser: true });

let Schema = mongoose.Schema;
let mongooseUniqueValidator = require("mongoose-unique-validator");

// Channel Model
let channelSchema = new Schema({
	name: {
        type: String,
        required: true,
        unique: true
    },
	description: {
        type: String,
        required: false
    },
	private: Boolean,
	members: [
		{
			_id: false,
			id: Number,
			username: String,
			firstname: String,
			lastname: String,
			photourl: String
		}
	],
	createdAt: {
        type: Date,
        default: Date.now
    },
	creator: {
		id: Number,
		username: String,
		firstname: String,
		lastname: String,
		photourl: String
	},
	editedAt: {
        type: Date
    }
});

channelSchema.plugin(mongooseUniqueValidator);
let ChannelData = mongoose.model("Channel", channelSchema);
channelSchema.set("toJSON", {
	transform: function (_, ret) {
		ret.id = ret._id;
		delete ret._id;
		delete ret.__v;
	}
});

// Message Model
app.use(express.json());
app.use(morgan("dev"));

let messageSchema = new Schema({
	_id: Number,
	channelID: ObjectId,
	body: {
        type: String,
        required: true
    },
	createdAt: {
        type: Date,
        default: Date.now
    },
	creator: {
		id: Number,
		username: String,
		firstname: String,
		lastname: String,
		photourl: String
	},
	editedAt: {
        type: Date
    }
});

messageSchema.plugin(autoIncrement);
messageSchema.plugin(mongooseUniqueValidator);
let Message = mongoose.model("Message", messageSchema);
messageSchema.set("toJSON", {
	transform: function (_, ret) {
		ret.id = ret._id;
		delete ret._id;
		delete ret.__v;
	}
});

// General Channel
let general = new ChannelData({
	name: "general",
	description: "This is the general channel."
});
general.save();

// APIs

// get all channels that the current user is allowed to see
app.get("/v1/channels", (req, res) => {
	let userJSON = req.get("X-User");
    if (userJSON) {
        let user = JSON.parse(userJSON);
        ChannelData.find({}).exec((err, channels) => {
            if (err) {
				res.statusCode = 400;
				return res.status(400).send("Error getting channels");
			}
            
            let channelView = [];
            for (let c of channels) {                
                if (!c.private || c.creator.id == user.id || isMember(c.members, user)) {
                    channelView.push(c);
                }
            }
            
			res.setHeader("Content-Type", "application/json");
			return res.status(200).send(JSON.stringify(channelView));
		});
	} else {
		return res.status(401).send("User not authenticated");
	}
});

// create new channel
app.post("/v1/channels", (req, res) => {
    // let userJSON = req.get("X-User");
	// if (userJSON) {
		// let user = JSON.parse(userJSON);
        let members = req.body.members;
        // members.push(user);
        
		if (req.body.name == undefined || req.body.name == "") {
            return res.status(400).send("Channel name is required");
        }
        
		let newChannel = new ChannelData({
			name: req.body.name,
			description: req.body.description,
			private: req.body.private || false,
			createdAt: new Date().toLocaleString(),
			// creator: user,
			members: members
        });
		newChannel.save(err => {
			if (err) {
				return res.status(400).send("Error saving new channel");
            }

            const newMsg = {
                "type": "channel-new",
                "channel": newChannel,
                "userIDs": newChannel.private ? getMemberID(newChannel.members) : []
            };
            channel.sendToQueue(
                queueName, 
                new Buffer(JSON.stringify(newMsg)),
                {persistent: true}
            );
			res.setHeader("Content-Type", "application/json");
			return res.status(201).send(JSON.stringify(newChannel));
		});
	// } else {
		// return res.status(401).send("User not authenticated");
	// }
});

// Specific channel
// GET
app.get("/v1/channels/:id", (req, res) => {
	let userJSON = req.get("X-User");
	if (userJSON) {
		let user = JSON.parse(userJSON);
		ChannelData.findOne({ "_id": req.params.id }).exec((err, channel) => {
			if (err || channel == undefined) {
                return res.status(400).send("Error finding channel");
            } else if (channel.private && !isMember(channel.members, user)) {
                return res.status(403).send("Not a member of this private channel");
            } else {
                let query = { channelID: req.params.id };
                if (req.params.before != undefined || req.params.before != "") {
                    query = { channelID: req.params.id, "_id": { $lt: req.params.before } };
                }
                Message.find(query).sort({ createdAt: -1 }).limit(100).exec((err, messages) => {
                    if (err) {
                        return res.status(500).send("Error finding messages")
                    } else {
                        res.set("Content-Type", "application/json");
                        return res.status(200).send(JSON.stringify(messages));
                    }
                });
            }
		});
	} else {
		return res.status(401).send("User not authenticated");
	}
});

// POST
app.post("/v1/channels/:id", (req, res) => {
    let userJSON = req.get("X-User");
	if (userJSON) {
		let user = JSON.parse(userJSON);
		ChannelData.findOne({ "_id": req.params.id }).exec((err, ch) => {
			if (err || ch == undefined) {
				return res.status(400).send("Error finding channel");
            } else if (ch.private && ch.creator.id != user.id && !isMember(ch.members, user)) {
                return res.status(403).send("Not a member of this private channel");
            } else {
                let msg = new Message({
                    channelID: new ObjectId(req.params.id),
                    body: req.body.body,
                    createdAt: new Date().toLocaleString(),
                    creator: user
                });
                msg.save(err => {
                    if (err) {
                        return res.status(500).send("Error saving new message");
                    } else {
                        const newMsg = {
                            "type": "message-new",
                            "message": msg,
                            "userIDs": ch.private ? getMemberID(ch.members) : []
                        };
                
                        channel.sendToQueue(
                            queueName, 
                            new Buffer(JSON.stringify(newMsg)),
                            {persistent: true}
                        );

                        res.setHeader("Content-Type", "application/json");
                        return res.status(201).send(JSON.stringify(msg));
                    }
                });
            }
		});
	} else {
		return res.status(401).send("User not authenticated");
	}
});

// PATCH
app.patch("/v1/channels/:id", (req, res) => {
	let userJSON = req.get("X-User");
	if (userJSON) {
		let user = JSON.parse(userJSON);
		ChannelData.findOne({ "_id": req.params.id }).exec((err, ch) => {
			if (err || ch == undefined) {
				res.setHeader("Content-Type", "text/plain");
				return res.status(400).send("Error finding channel")
			} else if (ch.creator.id != user.id) {
                res.setHeader("Content-Type", "text/plain");
                return res.status(403).send("The current user is not the creator of this channel");
            } else {
                ChannelData.update({ "_id": req.params.id }, { $set: { name: req.body.name || ch.name, description: req.body.description || ch.description } }, err => {
                    if (err) {
                        return res.status(500).send("Error updating")
                    } else {
                        ChannelData.findOne({ "_id": req.params.id }).exec((err, updated) => {
                            if (err) {
                                return res.status(500).send("Error updating")
                            } else {
                                const newMsg = {
                                    "type": "channel-update",
                                    "channel": updated,
                                    "userIDs": updated.private ? getMemberID(updated.members) : []
                                };
                        
                                channel.sendToQueue(
                                    queueName, 
                                    new Buffer(JSON.stringify(newMsg)),
                                    {persistent: true}
                                );
                                
                                res.setHeader("Content-Type", "application/json")
                                return res.status(200).send(JSON.stringify(updated))
                            }
                        });
                    }
                })
            }
		})
	} else {
		return res.status(401).send("User not authenticated");
	}
});

// DELETE
app.delete("/v1/channels/:id", (req, res) => {
	let userJSON = req.get("X-User");
	if (userJSON) {
		let user = JSON.parse(userJSON);
		ChannelData.findOne({ "_id": req.params.id }).exec((err, ch) => {
			if (err || ch == undefined) {
				res.setHeader("Content-Type", "text/plain");
				return res.status(400).send("Error finding channel")
			// } else if (ch.creator.id != user.id) {
            } else if (ch.members[0].id != user.id && ch.members[1].id != user.id) {
                res.setHeader("Content-Type", "text/plain");
                // return res.status(403).send("The current user is not the creator of this channel");
                return res.status(403).send("The current user is not a member of this channel");
            } else {
                Message.remove({ "channelID": req.params.id }, err => {
                    if (err) {
                        return res.status(500).send("Error removing message")
                    } else {
                        ChannelData.remove({ "_id": req.params.id }, err => {
                            if (err) {
                                return res.status(500).send("Error removing channel")
                            } else {
                                const newMsg = {
                                    "type": "channel-delete",
                                    "channelID": req.params.id,
                                    "userIDs": ch.private ? getMemberID(ch.members) : []
                                };
                        
                                channel.sendToQueue(
                                    queueName, 
                                    new Buffer(JSON.stringify(newMsg)),
                                    {persistent: true}
                                );

                                res.setHeader("Content-Type", "text/plain");
                                return res.status(200).send("Delete successful")
                            }
                        })
                    }
                })
            }
		})
	} else {
		return res.status(401).send("User not authenticated");
	}
});

// Members
// POST
app.post("/v1/channels/:id/members", (req, res) => {
	let userJSON = req.get("X-User");
	if (userJSON) {
		let user = JSON.parse(userJSON);
		ChannelData.findOne({ "_id": req.params.id }).exec((err, channel) => {
			if (err || channel == undefined) {
				res.setHeader("Content-Type", "text/plain");
				return res.status(400).send("Error finding channel")
			} else if (channel.creator.id != user.id) {
                res.setHeader("Content-Type", "text/plain");
                return res.status(403).send("The current user is not the creator of this channel");
            } else if (req.body.id == undefined || req.body.id == "") {
                res.setHeader("Content-Type", "text/plain");
                return res.status(400).send("User ID is required to be added as a member");
            } else {
                let members = channel.members
                members.push(req.body)
                ChannelData.update({ "_id": req.params.id }, { $set: { members: members } }, err => {
                    if (err) {
                        res.setHeader("Content-Type", "text/plain");
                        return res.status(500).send("Error adding member");
                    } else {
                        res.setHeader("Content-Type", "text/plain");
                        return res.status(201).send("User added as a member");
                    }
                });
            }
		})
	} else {
		return res.status(401).send("User not authenticated");
	}
});

// DELETE
app.delete("/v1/channels/:id/members", (req, res) => {
	let userJSON = req.get("X-User");
	if (userJSON) {
		let user = JSON.parse(userJSON);
		ChannelData.findOne({ "_id": req.params.id }).exec((err, channel) => {
			if (err || channel == undefined) {
				res.setHeader("Content-Type", "text/plain");
				return res.status(400).send("Error finding channel")
			} else if (channel.creator.id != user.id) {
                res.setHeader("Content-Type", "text/plain");
                return res.status(403).send("The current user is not the creator of this channel");
            } else if (req.body.id == undefined || req.body.id == "") {
                res.setHeader("Content-Type", "text/plain");
                return res.status(400).send("User ID of the member is required to be removed");
            } else {
                ChannelData.update({ "_id": req.params.id }, { $pull: { members: req.body } }, err => {
                    if (err) {
                        res.setHeader("Content-Type", "text/plain");
                        return res.status(500).send("Error deleting");
                    } else {
                        res.setHeader("Content-Type", "text/plain");
                        return res.status(200).send("Delete successful");
                    }
                });
            }
		})
	} else {
		return res.status(401).send("User not authenticated");
	}
});

// Specfic message
// PATCH
app.patch("/v1/messages/:messageID", (req, res) => {
	let userJSON = req.get("X-User");
	if (userJSON) {
		let user = JSON.parse(userJSON);
		Message.findOne({ "_id": req.params.messageID }).exec((err, message) => {
			if (err) {
				res.setHeader("Content-Type", "text/plain");
				return res.status(400).send("Error finding message")
			} else if (message.creator.id != user.id) {
                res.setHeader("Content-Type", "text/plain");
                return res.status(403).send("The current user is not the creator of this message");
            } else if (req.body.body == message.body) {
                res.setHeader("Content-Type", "text/plain");
                return res.status(400).send("No update");
            } else {
                Message.update({ "_id": req.params.messageID }, { $set: { body: req.body.body || message.body } }, err => {
                    if (err) {
                        return res.status(500).send("Error updating")
                    } else {
                        Message.findOne({ "_id": req.params.messageID }).exec((err, updated) => {
                            if (err) {
                                return res.status(500).send("Error updating")
                            } else {
                                ChannelData.findOne({ "_id": updated.channelID }).exec((err, ch) => {
                                    if (err || ch == undefined) {
                                        return res.status(400).send("Error finding channel");
                                    } else {
                                        const newMsg = {
                                            "type": "message-update",
                                            "message": updated,
                                            "userIDs": ch.private ? getMemberID(ch.members) : []
                                        };
                                
                                        channel.sendToQueue(
                                            queueName, 
                                            new Buffer(JSON.stringify(newMsg)),
                                            {persistent: true}
                                        );
                                    }
                                });
                                
                                res.setHeader("Content-Type", "application/json")
                                return res.status(200).send(JSON.stringify(updated))
                            }
                        });
                    }
                })
            }
		})
	} else {
		return res.status(401).send("User not authenticated");
	}
});

// DELETE
app.delete("/v1/messages/:messageID", (req, res) => {
	let userJSON = req.get("X-User");
	if (userJSON) {
		let user = JSON.parse(userJSON);
		Message.findOne({ "_id": req.params.messageID }).exec((err, message) => {
			if (err || message == undefined) {
				res.setHeader("Content-Type", "text/plain");
				return res.status(400).send("Error finding message")
			} else if (message.creator.id != user.id) {
                res.setHeader("Content-Type", "text/plain");
                return res.status(403).send("The current user is not the creator of this message");
            } else {
                Message.remove({ "_id": req.params.messageID }, err => {
                    if (err) {
                        return res.status(500).send("Error deleting")
                    } else {
                        ChannelData.findOne({ "_id": message.channelID }).exec((err, ch) => {
                            if (err || ch == undefined) {
                                return res.status(400).send("Error finding channel");
                            } else {
                                const newMsg = {
                                    "type": "message-delete",
                                    "messageID": req.params.messageID,
                                    "userIDs": ch.private ? getMemberID(ch.members) : []
                                };
                        
                                channel.sendToQueue(
                                    queueName, 
                                    new Buffer(JSON.stringify(newMsg)),
                                    {persistent: true}
                                );
                            }
                        });
                        
                        res.setHeader("Content-Type", "text/plain");
                        return res.status(200).send("Delete successful")
                    }
                })
            }
		})
	} else {
		return res.status(401).send("User not authenticated");
	}
});

app.use((err, _, res) => {
	console.error(err.stack)

	res.set("Content-Type", "text/plain");
	res.status(500).send(err.message);
});

// start the server listening on host:port
app.listen(portNum, host, () => {
	// callback is executed once server is listening
	console.log(`server is listening at http://${addr}...`);
});

function isMember(members, user) {
    for (let m of members) {
        if (m.id == user.id) {
            return true;
        }
    }
    return false;
}

function getMemberID(members) {
    let userID = [];
    for (let member of members) {
        userID.push(member.id);
    }
    return userID;
}
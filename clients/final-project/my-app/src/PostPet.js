import React from "react";
export default class PostPet extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            name:"",
            breed: "",
            gender: "",
            age:-1,
            bio:""
        };
    }
    componentDidMount() {
        let auth = localStorage.getItem("auth");
    }
    handleSignUp(e) {
        e.preventDefault();
        console.log(localStorage.getItem("auth"));
        fetch("https://api.demitu.me/v1/pet", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": localStorage.getItem("auth")
            },
            body: JSON.stringify({
                name: this.state.name,
                breed: this.state.breed.toLowerCase(),
                gender:this.state.gender,
                age:this.state.age,
                bio:this.state.bio
            }),
            mode: "cors",
            cache: "default"
        })
            .then(res => {
                if (!res.ok) {
                    throw Error(res.statusText + " " + res.status);
                }
                console.log(res.headers.get("Authorization"));
                return res.json();
            })
            .then(data => {
                console.log(data);
                localStorage.setItem("petName", data.name);
                localStorage.setItem("petAge", data.age);
                localStorage.setItem("petBio", data.bio);
                this.props.history.push({pathname: "/petlist"})
            })
            .catch(function(error) {
                alert(error);
            });
    }
    render(){
        return(
            <div>
                <header className="container-fluid text-white">
                    <div className="row ">
                        <div className="col-12 col-sm-12 col-md-12 col-lg-12 col-xl-12 pt-3 my-border">
                            <div className="text-center">
                                <h1>🔥Add Pet Profile🔥</h1>
                            </div>
                        </div>
                    </div>
                </header>
                <main>
                    <div className="d-flex justify-content-center pt-4 pb-5">
                        <div className="card w-75">
                            <div className="card-body">
                                <div className="container">
                                    <div>
                                        <div id="result" />
                                        <form className="form-group">
                                            <p>Name</p>
                                            {this.state.name === "" ? (
                                                <div className="alert alert-danger mt-2">
                                                    It shouldn't be blank
                                                </div>
                                            ) : (
                                                undefined
                                            )}
                                            <input
                                                id="Name"
                                                type="text"
                                                className="form-control"
                                                placeholder="Name"
                                                onInput={evt =>
                                                    this.setState({
                                                        name:
                                                            evt.target.value
                                                    })
                                                }
                                            />
                                        </form>
                                        <form className="form-group">
                                            <p>Breed</p>
                                            {this.state.breed === "" ? (
                                                <div className="alert alert-danger mt-2">
                                                    It shouldn't be blank
                                                </div>
                                            ) : (
                                                undefined
                                            )}
                                            <input
                                                id="Breed"
                                                type="text"
                                                className="form-control"
                                                placeholder="Breed"
                                                onInput={evt =>
                                                    this.setState({
                                                        breed:
                                                            evt.target.value
                                                    })
                                                }
                                            />
                                        </form>
                                        <form className="form-group">
                                            <p>Gender</p>
                                            {this.state.gender === "" ? (
                                                <div className="alert alert-danger mt-2">
                                                    It shouldn't be blank
                                                </div>
                                            ) : (
                                                undefined
                                            )}
                                            <input
                                                id="Gender"
                                                type="text"
                                                className="form-control"
                                                placeholder="Gender"
                                                onInput={evt =>
                                                    this.setState({
                                                        gender:
                                                            evt.target.value
                                                    })
                                                }
                                            />
                                        </form>
                                        <form className="form-group">
                                            <p>Age</p>
                                            {!Number.isInteger(this.state.age) ? (
                                                <div className="alert alert-danger mt-2">
                                                    Age should be a number
                                                </div>
                                            ) : (
                                                undefined
                                            )}
                                            <input
                                                id="Age"
                                                type="text"
                                                className="form-control"
                                                placeholder="Enter a number"
                                                onInput={evt =>
                                                    this.setState({
                                                        age:
                                                            parseInt(evt.target.value,10)
                                                    })
                                                }
                                            />
                                        </form>
                                        <form className="form-group">
                                            <p>Bio</p>
                                            <textarea
                                                id="Bio"
                                                type="text"
                                                className="form-control"
                                                placeholder="Bio"
                                                onInput={evt =>
                                                    this.setState({
                                                        bio:
                                                            evt.target.value
                                                    })
                                                }
                                            />
                                        </form>
                                        
                                        <button
                                            className="btn btn-primary mr-2 p-2"
                                            onClick={e => this.handleSignUp(e)}
                                        >Sign Up
                                        </button>
                                    </div>
                                    </div>
                            </div>
                        </div>
                    </div>
                </main>
            </div>
        )
    }
}
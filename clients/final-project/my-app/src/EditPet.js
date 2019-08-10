import React from "react";
import { toast, ToastContainer } from "react-toastify";
export default class EditPet extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            name:"",
            age: -1,
            bio:""
        }
    }
    handleSubmit(evt) {
        evt.preventDefault();
        fetch("https://api.demitu.me/v1/pet", {
            method: "PATCH",
            headers: {
                "Authorization": localStorage.getItem("auth"),
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                name: this.state.name,
                age: this.state.age,
                bio:this.state.bio
            })
        })
            .then(res => {
                if (!res.ok) {
                    throw Error(res.statusText + " " + res.status);
                }
                return res.json();
            })
            .then(data => {
                console.log(data);
                localStorage.setItem("petName", data.name);
                localStorage.setItem("petAge", data.age);
                localStorage.setItem("petBio", data.bio);
                toast("Update Success")
            })
            .catch(function(error) {
                alert(error);
            });
    }
    handleBack(evt){
        evt.preventDefault();
        this.props.history.push({pathname: "/getPet"});
    }
    render(){
        return(
            <div>
            <header className="container-fluid text-white">
                <div className="row">
                    <div className="col-12 col-sm-12 col-md-12 col-lg-12 col-xl-12 pt-3 my-border">
                        <div className="text-center">
                            <h1>🔥Edit Pet Profile🔥</h1>
                        </div>
                    </div>
                </div>
            </header>
            <ToastContainer/>
            <main>
                <div className="d-flex justify-content-center pt-4">
                    <div className="card w-75">
                        <div className="card-body">
                            <div className="container">
                                <div id="result" />
                                <form>
                                    <div className="form-group">
                                        <label htmlFor="Name">
                                            Name
                                        </label>
                                        <input
                                            type="text"
                                            id="name"
                                            className="form-control"
                                            placeholder={localStorage.getItem("petName")}
                                            onInput={evt =>
                                                this.setState({
                                                    name:
                                                        evt.target.value
                                                })
                                            }
                                        />
                                    </div>
                                    <div className="form-group">
                                        <label htmlFor="Age">
                                            Age
                                        </label>
                                        {!Number.isInteger(this.state.age) ? (
                                                <div className="alert alert-danger mt-2">
                                                    Age should be a number
                                                </div>
                                            ) : (
                                                undefined
                                        )}
                                        <input
                                            type="text"
                                            id="age"
                                            className="form-control"
                                            placeholder={localStorage.getItem("petAge")}
                                            onInput={evt =>
                                                this.setState({
                                                    age:
                                                        parseInt(evt.target.value,10)
                                                })
                                            }
                                        />
                                    </div>
                                    <div className="form-group">
                                        <label htmlFor="Name">
                                            Bio
                                        </label>
                                        <textarea
                                            type="text"
                                            id="name"
                                            className="form-control"
                                            placeholder={localStorage.getItem("petBio")}
                                            onInput={evt =>
                                                this.setState({
                                                    bio:
                                                        evt.target.value
                                                })
                                            }
                                        />
                                    </div>
                                    <div className="form-group">
                                        <button
                                            type="submit"
                                            className="btn btn-primary btn-block"
                                            onClick={evt =>
                                                this.handleSubmit(evt)
                                            }
                                        >
                                            Confirm
                                        </button>
                                    </div>
                                    <div className="form-group">
                                        <button
                                            type="submit"
                                            className="btn btn-secondary btn-block"
                                            onClick={evt =>
                                                this.handleBack(evt)
                                            }
                                        >
                                            Cancel & Back to Pet List
                                        </button>
                                    </div>
                                </form>
                            </div>
                        </div>
                    </div>
                </div>
            </main>
        </div>
        )
    }
}
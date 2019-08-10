import React from "react";
import { Link } from "react-router-dom";

export default class SignUp extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            email:"",
            password: "",
            passwordConfirm: "",
            userName: "",
            firstName: "",
            lastName: ""
        };
    }

    componentDidMount() {
        let auth = localStorage.getItem("auth");
        console.log(auth);
    }

    handleSignUp(e) {
        e.preventDefault();

        fetch("https://api.demitu.me/v1/users", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                email: this.state.email,
                password: this.state.password,
                passwordConf: this.state.passwordConfirm,
                userName: this.state.userName,
                firstName: this.state.firstName,
                lastName: this.state.lastName
            }),
            mode: "cors",
            cache: "default"
        })
            .then(res => {
                if (!res.ok) {
                    throw Error(res.statusText + " " + res.status);
                }
                console.log(res.headers.get("Authorization"));
                localStorage.setItem("auth", res.headers.get("Authorization"));
                return res.json();
            })
            .then(data => {
                console.log(data);
                localStorage.setItem("userid", data.id);
                localStorage.setItem("username", data.userName);
                this.props.history.push({pathname: "/getPet"})
            })
            .catch(function(error) {
                alert(error);
            });
    }

    render() {
        return (
            <div>
                <header className="container-fluid text-white">
                    <div className="row">
                        <div className="col-12 col-sm-12 col-md-12 col-lg-12 col-xl-12 pt-3 my-border">
                            <div className="text-center">
                                <h1>🔥Petinder🔥</h1>
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
                                            <p>Email</p>
                                            {this.state.email === "" ? (
                                                <div className="alert alert-danger mt-2">
                                                    It shouldn't be blank
                                                </div>
                                            ) : (
                                                undefined
                                            )}
                                            <input
                                                id="Email"
                                                type="text"
                                                className="form-control"
                                                placeholder="Email"
                                                onInput={evt =>
                                                    this.setState({
                                                        email:
                                                            evt.target.value
                                                    })
                                                }
                                            />
                                        </form>
                                        <form className="form-group">
                                            <p>User Name</p>
                                            {this.state.userName === "" ? (
                                                <div className="alert alert-danger mt-2">
                                                    It shouldn't be blank
                                                </div>
                                            ) : (
                                                undefined
                                            )}
                                            <input
                                                id="User Name"
                                                type="text"
                                                className="form-control"
                                                placeholder="User Name"
                                                onInput={evt =>
                                                    this.setState({
                                                        userName:
                                                            evt.target.value
                                                    })
                                                }
                                            />
                                        </form>
                                        <form className="form-group">
                                            <p>First Name</p>
                                            {this.state.firstName === "" ? (
                                                <div className="alert alert-danger mt-2">
                                                    It shouldn't be blank
                                                </div>
                                            ) : (
                                                undefined
                                            )}
                                            <input
                                                id="First Name"
                                                type="text"
                                                className="form-control"
                                                placeholder="First Name"
                                                onInput={evt =>
                                                    this.setState({
                                                        firstName:
                                                            evt.target.value
                                                    })
                                                }
                                            />
                                        </form>
                                        <form className="form-group">
                                            <p>Last Name</p>
                                            {this.state.lastName === "" ? (
                                                <div className="alert alert-danger mt-2">
                                                    It shouldn't be blank
                                                </div>
                                            ) : (
                                                undefined
                                            )}
                                            <input
                                                id="Last Name"
                                                type="text"
                                                className="form-control"
                                                placeholder="Last Name"
                                                onInput={evt =>
                                                    this.setState({
                                                        lastName:
                                                            evt.target.value
                                                    })
                                                }
                                            />
                                        </form>
                                        <form className="form-group">
                                            <p>Password</p>
                                            {this.state.password.length < 6 &&
                                            this.state.password.length !== 0 ? (
                                                <div className="alert alert-danger mt-2 text-wrap">
                                                    the password should be
                                                    longer than 6 digits
                                                </div>
                                            ) : (
                                                undefined
                                            )}
                                            <input
                                                id="password"
                                                type="password"
                                                className="form-control"
                                                placeholder="password"
                                                onInput={evt =>
                                                    this.setState({
                                                        password:
                                                            evt.target.value
                                                    })
                                                }
                                            />
                                        </form>
                                        <form className="form-group">
                                            <p>Password Confirmation</p>
                                            {this.state.password !==
                                                this.state.passwordConfirm &&
                                            this.state.passwordConfirm !==
                                                "" ? (
                                                <div className="alert alert-danger mt-2">
                                                    password should match
                                                </div>
                                            ) : (
                                                undefined
                                            )}
                                            <input
                                                id="passwordConfirm"
                                                type="password"
                                                className="form-control"
                                                placeholder="password"
                                                onInput={evt =>
                                                    this.setState({
                                                        passwordConfirm:
                                                            evt.target.value
                                                    })
                                                }
                                            />
                                        </form>
                                        <button
                                            className="btn btn-primary mr-2 p-2 btn-block"
                                            onClick={e => this.handleSignUp(e)}
                                        >
                                            Sign Up
                                        </button>
                                    </div>
                                </div>
                            </div>
                            <p className="text-center">Already have an account?{" "}<Link to="/login">Sign In!</Link></p>
                        </div>
                    </div>
                </main>
            </div>
        );
    }
}
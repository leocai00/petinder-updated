import React from "react";
import { Link } from "react-router-dom";
import { ROUTES } from "./constants";
export default class Login extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            email: "",
            password: ""
        };
    }

    componentDidMount() {
        let auth = localStorage.getItem("auth");
        console.log(auth);
    }

    handleSubmit(evt) {
        evt.preventDefault();
        fetch("https://api.demitu.me/v1/sessions", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                email: this.state.email,
                password: this.state.password
            })
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
                this.props.history.push({pathname: "/getPet"});

            })
            .catch(function(error) {
                alert(error);
            });
    }

    render() {
        return (
            <div>
                <header className="container-fluid text-white">
                    <div className="row ">
                        <div className="col-12 col-sm-12 col-md-12 col-lg-12 col-xl-12 pt-3 my-border">
                            <div className="text-center">
                                <h1>🔥Petinder🔥</h1>
                            </div>
                        </div>
                    </div>
                </header>
                <main>
                    <div className="d-flex justify-content-center pt-4">
                        <div className="card w-75">
                            <div className="card-body">
                                <div className="container">
                                    <div id="result" />
                                    <form>
                                        <div className="form-group">
                                            <label htmlFor="Email">
                                                Email
                                            </label>
                                            <input
                                                type="text"
                                                id="Email"
                                                className="form-control"
                                                placeholder="Email"
                                                onInput={evt =>
                                                    this.setState({
                                                        email:
                                                            evt.target.value
                                                    })
                                                }
                                                required
                                            />
                                        </div>
                                        <div className="form-group">
                                            <label htmlFor="password">
                                                Password
                                            </label>
                                            <input
                                                type="password"
                                                id="password"
                                                className="form-control"
                                                placeholder="Password"
                                                onInput={evt =>
                                                    this.setState({
                                                        password:
                                                            evt.target.value
                                                    })
                                                }
                                                required
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
                                                Sign In
                                            </button>
                                        </div>
                                    </form>
                                </div>
                            </div>
                            <p className="text-center">Don't have an account yet?{" "}<Link to={ROUTES.signUp}>{" "}Sign Up!</Link></p>
                        </div>
                    </div>
                </main>
            </div>
        );
    }
}
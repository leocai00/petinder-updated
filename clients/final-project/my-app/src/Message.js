import React from "react";
import ReactDOM from 'react-dom';
import {Button} from 'reactstrap';

import Mes from './Mes.js';
import { toast, ToastContainer } from "react-toastify";

export default class Message extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            socket: null,
            name: localStorage.getItem("username"),
            chats: []
        };

        this.submitMessage = this.submitMessage.bind(this);
        this.deleteChannel = this.deleteChannel.bind(this);
    }

    componentDidMount() {
        this.scrollToBot();
        let auth = localStorage.getItem("auth");
        let url = "wss://api.demitu.me/v1/ws?auth=" + auth;

        this.socket = new WebSocket(url);
        this.setState({ socket: this.socket });
        this.socket.onopen = () => {
            console.log("Connection Opened");
        };

        this.socket.onclose = () => {
            console.log("Connection Closed");
        };

        this.socket.onmessage = msg => {
            msg = JSON.parse(msg.data);
            if (msg.type === "channel-delete") {
                this.props.history.push({pathname: "/petlist"});
                toast("Chat ended!");
            } 
            if (msg.type === "message-new"){
                console.log("Message received " + msg.message.body);
                this.setState({
                    chats: this.state.chats.concat([{
                        username: msg.message.creator.id,
                        content: <p>{msg.message.body}</p>
                    }])
                }, () => {
                    ReactDOM.findDOMNode(this.refs.msg).value = "";
                });
            }
        }
    }

    componentDidUpdate() {
        this.scrollToBot();
    }

    scrollToBot() {
        ReactDOM.findDOMNode(this.refs.chats).scrollTop = ReactDOM.findDOMNode(this.refs.chats).scrollHeight;
    }

    submitMessage(e) {
        e.preventDefault();
        let channel = this.props.location.state.channel;
        console.log(channel)
        let url = "https://api.demitu.me/v1/channels/" + channel.id;
        fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": localStorage.getItem("auth")
            },
            body: JSON.stringify({
                body: ReactDOM.findDOMNode(this.refs.msg).value
            }),
            mode: "cors",
            cache: "default"
        })
            .then(res => {
                if (!res.ok) {
                    throw Error(res.statusText + " " + res.status);
                }
                return res.json();
            })
            .then(data => {
                console.log(data);
            })
            .catch(function(error) {
                alert(error);
            });
    }

    deleteChannel(e) {
        e.preventDefault();
        let channel = this.props.location.state.channel;
        console.log(channel)
        let url = "https://api.demitu.me/v1/channels/" + channel.id;
        fetch(url, {
            method: "DELETE",
            headers: {
                "Authorization": localStorage.getItem("auth")
            },
        })
        .then(res => {
            if (!res.ok) {
                throw Error(res.statusText + " " + res.status);
            }
            this.props.history.push({pathname: "/petlist"});
        })
        .catch(function(error) {
            alert(error);
        });
    }

    render() {
        const userID = localStorage.getItem("userid");
        return (
            <div className="App">
            <ToastContainer/>
            <Button className="exit" size="lg" onClick={evt =>this.deleteChannel(evt)}>X</Button>
                <div className="chatroom">
                    <h3 style={{'font-family':'Indie Flower'}}>🔥Petinder🔥</h3>
                    <ul className="chats" ref="chats">
                        {
                            this.state.chats.map((chat, index) => 
                                <Mes key={index} chat={chat} user={userID} />
                            )
                        }
                    </ul>
                    <form className="input" onSubmit={(e) => this.submitMessage(e)}>
                        <input type="text" ref="msg" />
                        <input type="submit" value="Submit" />
                    </form>
                </div>
            </div>
        );
    }
}


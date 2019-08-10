import React from "react";
import ReactCardFlip from 'react-card-flip';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
export class Msg extends React.Component {
    constructor(props) {
        super(props);
    }
    render(){
        return (
            <div>
                Hey! Someone you liked also liked you!
                <button onClick={evt => handleCreateChannel(evt, this.props.message, this.props.history)}>Start Conversation</button>
            </div>
        )
    }
}
function handleCreateChannel(evt, msg, history){
    evt.preventDefault()
    history.push({pathname: "/message", state: {channel: msg.channel}});
}

export default class PetList extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            pets:[],
            socket: null
        };
    }

    componentDidMount() {
        console.log(localStorage.getItem("auth"));
        fetch("https://api.demitu.me/v1/pet/matching", {
            method: "GET",
            headers: {
                "Authorization": localStorage.getItem("auth"),
            },
        })
        .then(res => {
            if (!res.ok) {
                throw Error(res.statusText + " " + res.status);
            }
            return res.json();
        })
        .then(data => {
            console.log(data);
            if (Object.keys(data).length !== 0){
                this.setState({pets: data})
            }
        })
        .catch(function(error) {
            alert(error);
        });

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
            if (msg.type === "Successfully matched") {
                toast(<Msg message={msg} history={this.props.history}/>, {autoClose: false})
            }

            if (msg.type === "pet-new") {
                let newPets = this.state.pets
                newPets.push(msg.pet)
                this.setState({pets: newPets})
            }
        };
    }
    handleSignOut(evt) {
        evt.preventDefault();
        fetch("https://api.demitu.me/v1/sessions/mine", {
            method: "DELETE",
            headers: {
                Authorization: localStorage.getItem("auth")
            }
        })
        .then(res => {
            if (!res.ok) {
                throw Error(res.statusText + " " + res.status);
            }
            localStorage.clear();
            this.props.history.push({ pathname: "/" });
        })
        .catch(function(error) {
            localStorage.clear();
        });
    }
    handleEditPet(evt){
        evt.preventDefault();
        this.props.history.push({pathname: "/editpet"});
    }
    render(){
        return(
            <div>
                <div className="clearfix" style={{ 'backgroundColor': '#ffbf84' }}>
                    <button className="btn btn-outline-secondary float-right" onClick={evt =>this.handleSignOut(evt)}>Sign out</button>
                    <button className="btn btn-outline-light float-right" onClick={evt =>this.handleEditPet(evt)}>Edit Pet Profile</button>
                </div>
                <List list={this.state.pets} />
                <ToastContainer/>
            </div>
        )
    }
}

export class List extends React.Component {
    constructor(props){
        super(props);
    }
    render(){
        let petList = this.props.list.map((p) => {
            return <PetCard key={p.name} info={p} adoptCallback={this.props.adoptCallback}/>
        })
        return(
            <div>
                <header className="container-fluid text-white">
                    <div className="row">
                        <div className="col-12 col-sm-12 col-md-12 col-lg-12 col-xl-12 pt-3 my-border">
                            <div className="text-center">
                                <h1>🔥Matching Pets🔥</h1>
                            </div>
                        </div>
                    </div>
                </header>
                <div className="card-deck p-4">
                    {petList}
                </div>
            </div>
        );
    }
}
export class PetCard extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            isFlipped: false
          };
        this.handleLike = this.handleLike.bind(this);
    }
    handleLike(evt, id){
        evt.preventDefault();
        this.setState(prevState => ({ isFlipped: !prevState.isFlipped }));
        let url = "https://api.demitu.me/v1/pet/" + id
        fetch(url, {
            method: "POST",
            headers: {
                "Authorization": localStorage.getItem("auth"),
            },
        })
        .then(res => {
            if (!res.ok) {
                throw Error(res.statusText + " " + res.status);
            }
            return res.json();
        })
        .then(data => {
            console.log(data);
            console.log("liked");
            localStorage.setItem("matchingid", id);
        })
        .catch(function(error) {
            alert(error);
        });
    }
    cancelLike(evt, id){
        evt.preventDefault();
        this.setState(prevState => ({ isFlipped: !prevState.isFlipped }));

        console.log(localStorage.getItem("auth"));
        let url = "https://api.demitu.me/v1/pet/" + id
        fetch(url, {
            method: "DELETE",
            headers: {
                "Authorization": localStorage.getItem("auth"),
            },
        })
        .then(res => {
            if (!res.ok) {
                throw Error(res.statusText + " " + res.status);
            }
        })
        .catch(function(error) {
            alert(error);
        });
    }
    render(){
        return(
            <ReactCardFlip isFlipped={this.state.isFlipped} flipDirection="horizontal">
                <div className="card" key="front">
                    <div className="card-body">
                        <h3 className="card-title">Name: {this.props.info.name}</h3>
                        <p className="card-text d-flex"> Gender: {this.props.info.gender}</p>
                        <p className="card-text d-flex"> Breed: {this.props.info.breed}</p>
                        <p className="card-text d-flex"> Age: {this.props.info.age}</p>
                        <p className="card-text d-flex"> Bio: {this.props.info.bio}</p>
                    </div>
                    <button style={{ 'backgroundColor': '#ffbf84' }} onClick={evt =>this.handleLike(evt, this.props.info.id)}>💖Click to like💖</button>
                </div>
                
    
                <div className="card" key="back">
                    <div className="card-body">
                        <h3 className="card-title">Successfully Liked!</h3>
                    </div>
                    <button style={{ 'backgroundColor': '#ffbf84' }} onClick={evt =>this.cancelLike(evt, this.props.info.id)}>✖️Click to Dislike✖️</button>
                </div>
            </ReactCardFlip>
        )
    }
}
import React from "react";
export default class GetPet extends React.Component {
    constructor(props) {
        super(props);
    }
    componentDidMount() {
        let auth = localStorage.getItem("auth");
        console.log(auth)
        fetch("https://api.demitu.me/v1/pet", {
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
                if (data === ""){
                    this.props.history.push({pathname: "/postpet"})
                } else {
                    localStorage.setItem("petName", data.name);
                    localStorage.setItem("petAge", data.age);
                    localStorage.setItem("petBio", data.bio);
                    this.props.history.push({pathname: "/petlist"})
                }
            })
            .catch(function(error) {
                alert(error);
            });
    }

    render(){
        return(
            <div></div>
        )
    }
}

import React from "react";
import { Jumbotron, Container, Button} from 'reactstrap';

export default class Main extends React.Component {
    constructor(props) {
        super(props);
    }

    handleSignin(evt){
        evt.preventDefault();
        this.props.history.push({ pathname: "/login"});
    }
    handleSignup(evt){
        evt.preventDefault();
        this.props.history.push({ pathname: "/signup"});
    }

    render(){
        return(
            <div>
                <Jumbotron className="text-center" fluid style={{ 'backgroundColor': '#ffbf84', height: 800}}>
                    <Container fluid>
                        <h1 className="display-4 text-center text-white">Say Hello to 🔥Petinder!🔥 </h1>
                        <p  className="text-center">🦍🙈🐶🐩🐷🐭🐰🐥🐍🦊🐱🦄🐟🐝🐔</p>
                        <p className="lead text-center text-white">Petinder is a social search web application that allows pet owners to like or dislike other pets, and allows pet owners to chat if both parties liked each other’s pet in the application. The app will be used as a pet's dating site.</p>
                        <hr className="my-2 text-center text-white" />
                        <p className="text-center text-white">Our small web application would be perfect for pet owners who want an easy, fun, yet serious way to find, not just a mate, but the perfect mate for their loving companion animals.</p>
                        <p className="lead text-center text-white">
                            <Button color="primary" size="lg" onClick={evt =>this.handleSignin(evt)}>Sign in</Button>
                            <Button color="secondary" size="lg" onClick={evt =>this.handleSignup(evt)}>Sign up</Button>
                        </p>
                    </Container>
                </Jumbotron>
            </div>
        )
    }
}

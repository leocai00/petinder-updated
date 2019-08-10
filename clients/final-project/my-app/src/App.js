import React, { Component } from 'react';
import './App.css';
import { BrowserRouter as Router, Route, Switch} from 'react-router-dom';
import Login from './Login';
import Signup from './Signup';
import { ROUTES } from "./constants";
import GetPet from './GetPet';
import PostPet from './PostPet';
import PetList from './PetList';
import Main from './Main';
import Message from './Message';
import EditPet from './EditPet';

class App extends Component {
  render() {
    return (
      <Router>
         <Switch>
            <Route exact path="/" component={Main} />
            <Route path={ROUTES.login} component={Login} />
            <Route path={ROUTES.signUp} component={Signup} />
            <Route path={ROUTES.getPet} component={GetPet} />
            <Route path={ROUTES.postPet} component={PostPet} />
            <Route path={ROUTES.petList} component={PetList} />
            <Route path={ROUTES.message} component={Message} />
            <Route path={ROUTES.editPet} component={EditPet} />
       </Switch>
      </Router>
    );
  }
}

export default App;

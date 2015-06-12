import React from 'react';
import AdminUserStore from '../stores/AdminUserStore';
import AdminUserActionCreator from '../actions/AdminUserActionCreators';
import App from './App.jsx';

export default React.createClass({
  _onChange() {
    this.setState(AdminUserStore.getAll());
  },

  getInitialState() {
    return AdminUserStore.getAll();
  },

  componentWillMount() {
    return AdminUserActionCreator.getUsers();
  },

  componentDidMount() {
    AdminUserStore.addChangeListener(this._onChange);
  },

  componentWillUnmount() {
    AdminUserStore.removeChangeListener(this._onChange);
  },

  render() {
    let {users} = this.state;
    return (
      <App users
        ={users} />
    );
  }
});

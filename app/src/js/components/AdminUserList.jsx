import React from 'react';
import AdminUser from './AdminUser.jsx';
import {Paper} from 'material-ui';

export default React.createClass({
  getDefaultProps() {
    return {
      users: []
    };
  },

  componentDidMount() {
  },

  render() {
    let {users} = this.props;
    return (
      <form id="user-list">
        {users.map(user =>
          <AdminUser user={user} />
        )}
      </form>
    );
  }
});

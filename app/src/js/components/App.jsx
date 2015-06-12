import React, {PropTypes} from 'react';
import AdminUserList from './AdminUserList.jsx';
import {Styles} from 'material-ui';

const ThemeManager = new Styles.ThemeManager();

export default React.createClass({
  propTypes: {
    users: PropTypes.array.isRequired
  },

  getDefaultProps() {
    return {
      users: []
    }
  },

  childContextTypes: {
    muiTheme: React.PropTypes.object
  },

  getChildContext() {
    return {
      muiTheme: ThemeManager.getCurrentTheme()
    };
  },

  render() {
    let {users} = this.props;
    return (
      <div className="example-page">
        <h1>Learning Flux</h1>
        <p>
          Below is a list of tasks you can implement to better grasp the patterns behind Flux.<br />
          Most features are left unimplemented with clues to guide you on the learning process.
        </p>

        <AdminUserList users={users} />
      </div>
    );
  }
});

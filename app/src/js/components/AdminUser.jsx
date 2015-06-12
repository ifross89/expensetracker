import React from 'react';
import {Menu} from 'material-ui'
export default React.createClass({
  getDefaultProps() {
    return {
      user: {
        name: '',
        email: '',
        active: false,
        admin: false
      }
    }
  },

  render() {
    let payload = [
      {
        payload: 1,
        text: 'Name',
        data: this.props.user.name
      },
      {
        payload: 2,
        text: 'Email',
        data: this.props.user.email
      }
    ];
    return (
      <Menu menuItems={payload} autoWidth={true} />
    );
  }
});

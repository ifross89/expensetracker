/* app.js: Entry point into the expensetracker app */

var React = require("react");
var _ = require("underscore");

var AdminUserStore = require("./stores/AdminUserStore");
var AdminUserActions = require("./actions/AdminUserActions");

var AdminUserBox = React.createClass({
	handleUserSubmit: function(user) {
		console.log("handleUserSubmit: " + user.email);
		AdminUserActions.saveUser(user);
	},

	handleUserDelete: function(userId) {
		console.log("Deleting user with id = " + userId);
		AdminUserActions.deleteUser(userId);
	},

	getInitialState: function() {
		return {data: AdminUserStore.getUsers()};
	},

	componentDidMount: function() {
		AdminUserStore.addChangeListener(this._onChange);
		AdminUserActions.getUsers();
		setInterval(AdminUserActions.getUsers, this.props.pollInterval);
	},

	componentWillUnmount: function() {
		AdminUserStore.removeChangeListener(this._onChange);
	},

	_onChange: function() {this.setState({data: AdminUserStore.getUsers()});},

	render: function() {
		return(
			<div className="adminUserBox">
				<h2> Users </h2>
				<UserList data={this.state.data} handleDelete={this.handleUserDelete}/>
				<NewUserAdminForm onUserSubmit={this.handleUserSubmit} />
			</div>
		)
	}
});

function asArr(obj) {
	var arr = [];
	_.each(obj, function(prop) {
		arr = arr.concat(prop);
	})
	return arr
}

var UserList = React.createClass({
	handleDelete: function(userId) {
		this.props.handleDelete(userId);
	},
	render: function() {
			var handleDelete = this.props.handleDelete;
			var userNodes = asArr(this.props.data).map(function(user){
				var userId = user.id ? user.id : -1;
				return (
					<User name={user.name} email={user.email} key={user.key} userId={userId} handleDelete={handleDelete}/>
				);
			});

			return (
			<div className="userList">
				<h3>Hello world, I am an UserList</h3>
				<ul>
					{userNodes}
				</ul>
			</div>
		);
	}
});

var User = React.createClass({
	handleDelete: function(e) {
		e.preventDefault();
		this.props.handleDelete(this.props.userId);
	},
	render: function() {
		var email = 'mailto:' + this.props.email;
		var deleteDisabled = this.props.userId == -1;
		return (
			<li className="user" ><a href={email}>{this.props.name}</a> <button disabled={deleteDisabled} onClick={this.handleDelete}>Delete</button></li>
		);
	}
});

var NewUserAdminForm = React.createClass({
	handleSubmit: function(e) {
		e.preventDefault();
		var name = this.refs.name.getDOMNode().value.trim();
		var email = this.refs.email.getDOMNode().value.trim();
		var admin = this.refs.admin.getDOMNode().value.trim() === "on";
		var active = this.refs.active.getDOMNode().value.trim() === "on";
		var password = this.refs.password.getDOMNode().value.trim();

		if (!name || !email) {
			return;
		}


		console.log("handleSubmit: " + email);
		this.props.onUserSubmit({
			name: name,
			email: email,
			admin: admin,
			active: active,
			password: password,
		});

		this.refs.name.getDOMNode().value = '';
		this.refs.email.getDOMNode().value = '';
		this.refs.admin.getDOMNode().value = 'off';
		this.refs.active.getDOMNode().value = 'off';
		this.refs.password.getDOMNode().value = '';
		return;
	},
	render: function() {
		return (
			<form className="newUserAdminForm" onSubmit={this.handleSubmit}>
				<input type="text" placeholder="Your name" name="name" ref="name"/><br />
				<input type="email" placeholder="hello@example.com" name="email" ref="email" /><br />
				<input type="password" name="password" ref="password" /> <br />
				<input type="checkbox" name="admin" ref="admin"/> Admin? <br />
				<input type="checkbox" name="active" ref="active" /> Active? <br />
				<input type="submit" value="Save" /><br />
			</form>
		);
	}
});

React.render(
	<AdminUserBox url="/admin/users" pollInterval={1000000} />,
	document.getElementById('content')
);
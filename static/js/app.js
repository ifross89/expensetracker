/* app.js: Entry point into the expensetracker app */

var React = require("react");
var _ = require("underscore");

var AdminUserStore = require("./stores/admin-user-store");
var AdminUserActions = require("./actions/admin-user-actions");
var AuthActions = require('./actions/auth-actions');

var SigninForm = React.createClass({
	handleSubmit: function(e) {
		e.preventDefault();

		var email = this.refs.email.getDOMNode().value.trim();
		var password = this.refs.password.getDOMNode().value.trim();

		console.log('handle login submit', email, password);
		AuthActions.login(email, password);
	},

	render: function() {
		return (
			<div>
				<form onSubmit={this.handleSubmit}>
					<input type="email" name="email" ref="email" />
					<input type="password" name="password" ref="password" />
					<input type="submit" value="Log In" /><br />
				</form>
			</div>
		)
	}
});

var LogoutButton = React.createClass({
	handleSubmit: function(e) {
		e.preventDefault();

		AuthActions.logout();
	},

	render: function() {
		return (
			<div>
				<form onSubmit={this.handleSubmit}>
					<input type="submit" value="Log Out" />
				</form>
			</div>
		)
	}
});

var ChangePasswordForm = React.createClass({
	handleSubmit: function(e) {
		e.preventDefault();

		var currentPassword = this.refs.oldPassword.getDOMNode().value.trim();
		var newPassword = this.refs.newPassword.getDOMNode().value.trim();
		var repeatPassword = this.refs.confirmPassword.getDOMNode().value.trim();
		AuthActions.changePassword(
			currentPassword, newPassword, repeatPassword
		);
	},

	render: function() {
		return (
			<div>
				<form onSubmit={this.handleSubmit}>
					<input type="password" name="oldPassword" ref="oldPassword" />
					<input type="password" name="newPassword" ref="newPassword" />
					<input type="password" name="confirmPassword" ref="confirmPassword" />
					<input type="submit" value="Change password" />
					</form>
					</div>
		)
	}
});

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
				<SigninForm />
				<LogoutButton />
				<ChangePasswordForm />
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
	});
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
	getInitialState: function() {
		return {
			active: false,
			admin: false
		}
	},
	handleActiveChanged: function(e) {
		var state = this.state;
		state.active = !state.active;

		this.setState(state);
	},
	handleAdminChanged: function() {
		var state = this.state;
		state.admin = !state.admin

		this.setState(state);
	},
	handleSubmit: function(e) {
		e.preventDefault();
		var name = this.refs.name.getDOMNode().value.trim();
		var email = this.refs.email.getDOMNode().value.trim();
		var admin = this.state.admin;
		var active = this.state.active;
		var password = this.refs.password.getDOMNode().value.trim();



		if (!name || !email) {
			return;
		}


		this.props.onUserSubmit({
			name: name,
			email: email,
			admin: admin,
			active: active,
			password: password
		});

		this.refs.name.getDOMNode().value = '';
		this.refs.email.getDOMNode().value = '';
		this.refs.password.getDOMNode().value = '';
		this.setState({active:false, admin:false});
	},
	render: function() {
		return (
			<form className="newUserAdminForm" onSubmit={this.handleSubmit}>
				<input type="text" placeholder="Your name" name="name" ref="name"/><br />
				<input type="email" placeholder="hello@example.com" name="email" ref="email" /><br />
				<input type="password" name="password" ref="password" /> <br />
				<input type="checkbox" name="admin" ref="admin" onChange={this.handleAdminChanged} checked={this.state.admin}/> Admin? <br />
				<input type="checkbox" name="active" ref="active" onChange={this.handleActiveChanged} checked={this.state.active} /> Active? <br />
				<input type="submit" value="Save" /><br />
			</form>
		);
	}
});

React.render(
	<AdminUserBox url="/admin/users" pollInterval={1000000} />,
	document.getElementById('content')
);

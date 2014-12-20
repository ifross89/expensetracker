var _ = require('underscore');
var $ = require('jquery');
var React = require('react');

var initialData = [
	{
		id: 1,
		name: "Warehouse",
		users:
		[
			{
				id: 1,
				email: "ifross@gmail.com",
				name: "Ian",
			},
			{
				id: 2,
				email: "giles@hello.com",
				name: "Giles Moody",
			},
		],
		expenses:
		[
			{
				id: 1,
				amount: 100,
				payerId: 1,
				groupId: 1,
				category: "Alcohol",
				description: "TINNEHS",
				createdAt: "2012-04-23T18:25:43.511Z",
				assignments:
				[
					{
						id: 1,
						userId: 1,
						amount: 50,
						expenseId: 1,
						groupId: 1,
					},
					{
						id: 2,
						userId: 2,
						amount: 50,
						expenseId: 1,
						groupId: 1,
					}
				],
			},
			{
				id: 2,
				amount: 200,
				payerId: 2,
				groupId: 1,
				category: "Groceries",
				description: "Steak",
				createdAt: "2012-04-24T18:25:43.511Z",
				assignments:
				[
					{
						id: 3,
						userId: 1,
						expenseId: 2,
						groupId: 1,
						amount: 100,
					},
					{
						id: 4,
						userId: 2,
						expenseId: 2,
						groupId: 1,
						amount: 100,
					},
				]
			}
		],
		payments:
		[
			{
				id: 1,
				groupId: 1,
				amount: 20,
				giverId: 2,
				receiverId: 1,
				createdAt: "2012-04-26T18:25:43.511Z",
			},
			{
				id: 2,
				groupId: 1,
				amount: 40,
				giverId: 1,
				receiverId: 2,
				createdAt: "2012-04-26T18:25:43.511Z",
			}
		]
	},
];

function Keyer(prefix) {
	var prefix = prefix + "-";
	var i = 0;
	this.array = function(arr) {
		_.each(arr, function(obj){
			obj.key=prefix + obj.id;
		});
	}

	this.pending = function(obj) {
		obj.key = "pending-" + prefix + this.i;
		i++;
	}
}

var ExpenseTracker = React.createClass({
	getInitialState: function() {
		return {data: initialData};
	},
	render: function() {
		return (
			<div className="epenseTracker">
				<h1> ExpenseTracker </h1>
				<GroupList data={this.state.data} />
				<NewGroupForm />
			</div>
		);
	}
});

var GroupList = React.createClass({
	render: function() {
		var groupNodes = this.props.data.map(function(group) {
			return (
				<GroupEntry data={group} />
			);
		});
		return (
			<div className="groupList">
				{groupNodes}
			</div>
		);
	}
});

var GroupEntry = React.createClass({
	render: function() {
		return (
			<div className="groupEntry">
				Hello world, I am a GroupEntry
				<UserList data={this.props.data.users} />
			</div>
		);
	}
});

var UserList = React.createClass({
	render: function() {
			var userNodes = this.props.data.map(function(user){
				return (
					<User name={user.name} email={user.email} key={user.key} />
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
	render: function() {
		var email = 'mailto:' + this.props.email;
		return (
			<li className="user" ><a href={email}>{this.props.name}</a></li>
		);
	}
});

var AdminUserBox = React.createClass({
	keyer: new Keyer("user"),
	loadUsersFromServer: function() {
		$.ajax({
			url: this.props.url,
			dataType: 'json',
			success: function(response) {
				this.keyer.array(response.data);
				this.setState({data: response.data});
			}.bind(this),
			error: function(xhr, status, err) {
				console.error(this.props.url, status, err.toString());
			}.bind(this)
		});
	},
	handleUserSubmit: function(user) {
		var users = this.state.data;
		this.keyer.pending(user);
		users.push(user);
		this.setState({data: users}, function(){
			// setState accepts a callback. To avoid (improbable) race condition,
			// we'll send the ajax request right after we optimistically set the new
			// state
			$.ajax({
				url: this.props.url,
				dataType: 'json',
				type: 'POST',
				data: JSON.stringify(user),
				success: function(response) {
					this.keyer.array(response.data);
					this.setState({data: response.data});
				}.bind(this),
				error: function(xhr, status, err) {
					console.error(this.props.url, status, err.toString());
				}.bind(this)
			});
		});
	},
	getInitialState: function() {
		return {data: []};
	},
	componentDidMount: function() {
		this.loadUsersFromServer();
		setInterval(this.loadUsersFromServer, this.props.pollInterval);
	},
	render: function() {
		return(
			<div className="adminUserBox">
				<h2> Users </h2>
				<UserList data={this.state.data} />
				<NewUserAdminForm onUserSubmit={this.handleUserSubmit} />
			</div>
		)
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



var NewGroupForm = React.createClass({
	render: function() {
		return (
			<div className="newGroupForm">
				Hello world! I am a NewGroupForm
			</div>
		);
	}
});

React.render(
	<AdminUserBox url="/admin/users" pollInterval={100000} />,
	document.getElementById('content')
);
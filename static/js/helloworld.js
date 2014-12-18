var ExpenseTracker = React.createClass({
	render: function() {
		return (
			<div className="epenseTracker">
				<h1> ExpenseTracker </h1>
				<GroupList />
				<NewGroupForm />
			</div>
		);
	}
});

var GroupList = React.createClass({
	render: function() {
		return (
			<div className="groupList">
				Hello world! I am a GroupList
				<GroupEntry />
			</div>
		);
	}
});

var GroupEntry = React.createClass({
	render: function() {
		return (
			<div className="groupEntry">
				Hello world, I am a GroupEntry
				<UserList />
			</div>
		);
	}
});

var UserList = React.createClass({
	render: function() {
		return (
			<div className="userList">
				<h3>Hello world, I am an UserList</h3>
				<ul>
					<User email="ifross@gmail.com" />
					<User email="hello@example.com" />
				</ul>
			</div>
		);
	}
});

var User = React.createClass({
	render: function() {
		return (
			<li className="user" >{this.props.email}</li>
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
	<ExpenseTracker />, 
	document.getElementById('content')
);
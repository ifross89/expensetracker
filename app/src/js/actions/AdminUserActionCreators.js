import Dispatcher from '../Dispatcher';
import Constants from '../Constants';
import AdminUserClient from '../clients/AdminUserClient';
import utils from './ActionCreatorUtils';

export default {
  getUsers: () => {
    Dispatcher.handleViewAction({
      type: Constants.ActionTypes.GET_ADMIN_USERS
    });

    AdminUserClient.getAll()
      .then(utils.createServerActionHandler(Constants.ActionTypes.GET_ADMIN_USERS_SUCCESS),
            utils.createServerActionHandler(Constants.ActionTypes.GET_ADMIN_USERS_FAIL))
  },

  clearList() {
    console.warn('clearList action not yet implemented...');
  },

  completeTask(task) {
    console.warn('completeTask action not yet implemented...');
  }
};

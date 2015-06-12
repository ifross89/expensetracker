'use strict';

import Dispatcher from '../Dispatcher';

export default {
  createServerActionHandler(actionType) {
    return payload => {
      Dispatcher.handleServerAction({
        type: actionType,
        payload: payload
      });
    };
  }
}

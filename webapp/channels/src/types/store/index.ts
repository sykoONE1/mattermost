// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {GlobalState as BaseGlobalState} from '@mattermost/types/store';

import type * as MMReduxTypes from 'mattermost-redux/types/actions';

import type {PluginsState} from './plugins';
import type {ViewsState} from './views';

export type DraggingState = {
    state?: string;
    type?: string;
    id?: string;
}

export type GlobalState = BaseGlobalState & {
    plugins: PluginsState;
    storage: {
        storage: Record<string, any>;
        initialized: boolean;
    };
    views: ViewsState;
};

/**
 * A version of {@link MMReduxTypes.DispatchFunc} which supports dispatching web app actions.
 */
export type DispatchFunc = MMReduxTypes.DispatchFunc;

/**
 * A version of {@link MMReduxTypes.GetStateFunc} which supports web app state.
 */
export type GetStateFunc<State extends GlobalState = GlobalState> = MMReduxTypes.GetStateFunc<State>;

/**
 * A version of {@link MMReduxTypes.ActionFunc} which supports web app state and allows dispatching its actions.
 */
export type ActionFunc<
    Data = unknown,
    State extends GlobalState = GlobalState,
> = MMReduxTypes.ActionFunc<Data, State>;

/**
 * A version of {@link MMReduxTypes.ActionFuncAsync} which supports web app state and allows dispatching its actions.
 */
export type ActionFuncAsync<
    Data = unknown,
    State extends GlobalState = GlobalState,
> = MMReduxTypes.ActionFuncAsync<Data, State>;

/**
 * A version of {@link MMReduxTypes.ThunkActionFunc} which supports web app state and allows dispatching its actions.
 */
export type ThunkActionFunc<
    ReturnType,
    State extends GlobalState = GlobalState
> = MMReduxTypes.ThunkActionFunc<ReturnType, State>;

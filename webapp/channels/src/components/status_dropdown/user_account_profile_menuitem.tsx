// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useMemo} from 'react';
import {FormattedMessage} from 'react-intl';
import {useDispatch, useSelector} from 'react-redux';

import {AccountOutlineIcon} from '@mattermost/compass-icons/components';
import type {UserProfile} from '@mattermost/types/users';

import {savePreferences} from 'mattermost-redux/actions/preferences';
import {getInt} from 'mattermost-redux/selectors/entities/preferences';

import {openModal} from 'actions/views/modals';

import * as Menu from 'components/menu';
import {OnboardingTaskCategory, OnboardingTasksName, TaskNameMapToSteps, CompleteYourProfileTour} from 'components/onboarding_tasks';
import UserSettingsModal from 'components/user_settings/modal';

import {ModalIdentifiers} from 'utils/constants';

import type {GlobalState} from 'types/store';

interface Props {
    userId: UserProfile['id'];
}

export default function UserAccountProfileMenuItem(props: Props) {
    const dispatch = useDispatch();

    const onboardingTaskStep = useSelector((state: GlobalState) => getInt(state, OnboardingTaskCategory, OnboardingTasksName.COMPLETE_YOUR_PROFILE, 0));
    const isCompleteYourProfileTaskPending = onboardingTaskStep === TaskNameMapToSteps[OnboardingTasksName.COMPLETE_YOUR_PROFILE].STARTED;

    function handleClick() {
        dispatch(openModal({
            modalId: ModalIdentifiers.USER_SETTINGS,
            dialogType: UserSettingsModal,
            dialogProps: {isContentProductSettings: false},
        }));
    }

    function handleTourClick() {
        const taskName = OnboardingTasksName.COMPLETE_YOUR_PROFILE;
        const steps = TaskNameMapToSteps[taskName];

        dispatch(savePreferences(props.userId, [{
            user_id: props.userId,
            category: OnboardingTaskCategory,
            name: taskName,
            value: steps.FINISHED.toString(),
        }]));
    }

    const trailingElement = useMemo(() => {
        if (isCompleteYourProfileTaskPending) {
            return (
                <div onClick={handleTourClick}>
                    <CompleteYourProfileTour/>
                </div>
            );
        }

        return null;
    }, [isCompleteYourProfileTaskPending]);

    return (
        <>
            <Menu.Item
                id='userAccountPopover.menuItem.profile'
                leadingElement={
                    <AccountOutlineIcon
                        size={18}
                    />
                }
                labels={
                    <FormattedMessage
                        id='userAccountPopover.menuItem.profile'
                        defaultMessage='Profile'
                    />
                }
                trailingElements={trailingElement}
                onClick={handleClick}
            />
            <Menu.Separator/>
        </>
    );
}
// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {useEffect, useState} from 'react';

import type {DesktopNotificationPermission} from 'components/common/hooks/use_desktop_notification_permission';
import {getDesktopAppNotificationPermission} from 'components/common/hooks/use_desktop_notification_permission';
import NotificationPermissionDeniedNotice from 'components/user_settings/notifications/desktop_and_mobile_notification_setting/notification_permission_section_notice/notification_permission_denied_section_notice';
import NotificationPermissionNeverGrantedNotice from 'components/user_settings/notifications/desktop_and_mobile_notification_setting/notification_permission_section_notice/notification_permission_never_granted_section_notice';
import NotificationPermissionUnsupportedSectionNotice from 'components/user_settings/notifications/desktop_and_mobile_notification_setting/notification_permission_section_notice/notification_permission_unsupported_section_notice';

import {getNotificationPermission, isNotificationAPISupported, NotificationPermissionDenied, NotificationPermissionNeverGranted} from 'utils/notifications';

import NotificationPermissionDesktopDeniedSectionNotice from './notification_permission_desktop_denied_section_notice';

export default function NotificationPermissionSectionNotice() {
    const isNotificationSupported = isNotificationAPISupported();

    const [notificationPermission, setNotificationPermission] = useState(getNotificationPermission());
    const [desktopNotificationPermission, setDesktopNotificationPermission] = useState<DesktopNotificationPermission>(undefined);

    useEffect(() => {
        async function getDesktopAppNotificationPermissionAndSetState() {
            const permission = await getDesktopAppNotificationPermission();
            setDesktopNotificationPermission(permission);
        }

        getDesktopAppNotificationPermissionAndSetState();
    }, []);

    function handleRequestNotificationClicked(permission: NotificationPermission) {
        setNotificationPermission(permission);
    }

    if (!isNotificationSupported) {
        return <NotificationPermissionUnsupportedSectionNotice/>;
    }

    if (desktopNotificationPermission === NotificationPermissionDenied) {
        return <NotificationPermissionDesktopDeniedSectionNotice/>;
    }

    if (isNotificationSupported && notificationPermission === NotificationPermissionNeverGranted) {
        return <NotificationPermissionNeverGrantedNotice onCtaButtonClick={handleRequestNotificationClicked}/>;
    }

    if (isNotificationSupported && notificationPermission === NotificationPermissionDenied) {
        return <NotificationPermissionDeniedNotice/>;
    }

    return null;
}


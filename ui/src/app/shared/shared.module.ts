import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {NotFoundComponent} from './not-found/not-found.component';
import {CoreModule} from '../core/core.module';
import {AuthUserActiveService} from './route/auth-user-active.service';
import {NullFilterPipe} from './pipe/null-filter.pipe';
import {DeleteAlertComponent} from './common-component/delete-alert/delete-alert.component';
import {ConfirmAlertComponent} from './common-component/confirm-alert/confirm-alert.component';
import {StatusPipe} from './pipe/status.pipe';
import {StatusColorPipe} from './pipe/status-color.pipe';
import {ModalAlertComponent} from './common-component/modal-alert/modal-alert.component';
import {KeysPipe} from './pipe/keys.pipe';
import {ProgressColorPipe} from './pipe/progress-color.pipe';
import {BadgeClassPipe} from './pipe/badge-class.pipe';
import {DonePipe} from './pipe/done.pipe';
import {MessageTypePipe} from './pipe/message-type.pipe';
import {SubscribeCheckPipe} from './pipe/subscribe-check.pipe';
import {ReadStatusPipe} from './pipe/read-status.pipe';
import {MessageLevelPipe} from './pipe/message-level.pipe';
import {UploadStatusPipe} from './pipe/upload-status.pipe';
import {UploadComponent} from './common-component/upload/upload.component';
import { MessageDetailPipe } from './pipe/message-detail.pipe';

@NgModule({
  declarations: [NotFoundComponent, KeysPipe, NullFilterPipe, DeleteAlertComponent,
    ConfirmAlertComponent, StatusPipe, StatusColorPipe, ModalAlertComponent, ProgressColorPipe,
    BadgeClassPipe, DonePipe, MessageTypePipe, SubscribeCheckPipe, ReadStatusPipe, MessageLevelPipe, UploadComponent, UploadStatusPipe, MessageDetailPipe],
  imports: [
    CommonModule,
    CoreModule
  ], exports: [
        CoreModule,
        NullFilterPipe,
        StatusPipe,
        DeleteAlertComponent,
        ConfirmAlertComponent,
        ModalAlertComponent,
        StatusColorPipe,
        KeysPipe,
        ProgressColorPipe,
        BadgeClassPipe,
        DonePipe,
        MessageTypePipe,
        SubscribeCheckPipe,
        ReadStatusPipe,
        MessageLevelPipe,
        UploadComponent,
        UploadStatusPipe,
        MessageDetailPipe
    ], providers: [
    AuthUserActiveService,
  ]
})
export class SharedModule {
}

import {NgModule} from '@angular/core';
import {MessageCenterComponent} from './message-center.component';
import {RouterModule} from '@angular/router';
import {CoreModule} from '../../core/core.module';
import {UserReceiverComponent} from './user-receiver/user-receiver.component';
import {UserSubscribeComponent} from './user-subscribe/user-subscribe.component';
import {SharedModule} from '../../shared/shared.module';
import {MailboxComponent} from './mailbox/mailbox.component';
import {MailboxListComponent} from './mailbox/mailbox-list/mailbox-list.component';
import {MailboxDetailComponent} from './mailbox/mailbox-detail/mailbox-detail.component';
import {MailboxDeleteComponent} from './mailbox/mailbox-delete/mailbox-delete.component';


@NgModule({
    declarations: [MessageCenterComponent, UserReceiverComponent, UserSubscribeComponent,
        MailboxComponent, MailboxListComponent, MailboxDetailComponent, MailboxDeleteComponent],
    imports: [
        RouterModule,
        CoreModule,
        SharedModule
    ]
})
export class MessageCenterModule {
}

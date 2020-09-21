import {NgModule} from '@angular/core';
import {MessageCenterComponent} from './message-center.component';
import {RouterModule} from '@angular/router';
import {CoreModule} from '../../core/core.module';
import {UserReceiverComponent} from './user-receiver/user-receiver.component';
import { UserSubscribeComponent } from './user-subscribe/user-subscribe.component';
import {SharedModule} from '../../shared/shared.module';


@NgModule({
    declarations: [MessageCenterComponent, UserReceiverComponent, UserSubscribeComponent],
    imports: [
        RouterModule,
        CoreModule,
        SharedModule
    ]
})
export class MessageCenterModule {
}

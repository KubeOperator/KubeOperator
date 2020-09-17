import {NgModule} from '@angular/core';
import {MessageCenterComponent} from './message-center.component';
import {RouterModule} from '@angular/router';
import {CoreModule} from '../../core/core.module';
import {UserReceiverComponent} from './user-receiver/user-receiver.component';


@NgModule({
    declarations: [MessageCenterComponent, UserReceiverComponent],
    imports: [
        RouterModule,
        CoreModule
    ]
})
export class MessageCenterModule {
}

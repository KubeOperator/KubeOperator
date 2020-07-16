import {NgModule} from '@angular/core';
import {SettingComponent} from './setting.component';
import {RouterModule} from '@angular/router';
import {CoreModule} from '../../core/core.module';
import {CredentialModule} from './credential/credential.module';
import { SystemComponent } from './system/system.component';


@NgModule({
    declarations: [SettingComponent, SystemComponent],
    imports: [
        RouterModule,
        CoreModule,
        CredentialModule,
    ]
})

export class SettingModule {

}

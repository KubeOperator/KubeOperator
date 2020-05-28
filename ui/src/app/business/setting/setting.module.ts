import {NgModule} from '@angular/core';
import {SettingComponent} from './setting.component';
import {RouterModule} from '@angular/router';
import {CoreModule} from '../../core/core.module';
import {CredentialModule} from './credential/credential.module';


@NgModule({
    declarations: [SettingComponent],
    imports: [
        RouterModule,
        CoreModule,
        CredentialModule,
    ]
})

export class SettingModule {

}

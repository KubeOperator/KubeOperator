import {NgModule} from '@angular/core';
import {SettingComponent} from './setting.component';
import {RouterModule} from '@angular/router';
import {CoreModule} from '../../core/core.module';
import {CredentialModule} from './credential/credential.module';
import { SystemComponent } from './system/system.component';
import {BackupAccountModule} from './backup-account/backup-account.module';
import { LicenseComponent } from './license/license.component';


@NgModule({
    declarations: [SettingComponent, SystemComponent, LicenseComponent],
    imports: [
        RouterModule,
        CoreModule,
        CredentialModule,
        BackupAccountModule,
    ]
})

export class SettingModule {

}

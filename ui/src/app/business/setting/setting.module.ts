import {NgModule} from '@angular/core';
import {SettingComponent} from './setting.component';
import {RouterModule} from '@angular/router';
import {CoreModule} from '../../core/core.module';
import {CredentialModule} from './credential/credential.module';
import {SystemComponent} from './system/system.component';
import {BackupAccountModule} from './backup-account/backup-account.module';
import {LicenseComponent} from './license/license.component';
import {LicenseImportComponent} from './license/license-import/license-import.component';
import {LdapComponent} from './ldap/ldap.component';
import {ThemeComponent} from './theme/theme.component';
import {MessageComponent} from './message/message.component';
import {EmailComponent} from './email/email.component';
import { RegistrySettingComponent } from './registry-setting/registry-setting.component';


@NgModule({
    declarations: [SettingComponent, SystemComponent, LicenseComponent, LicenseImportComponent,
        LdapComponent, ThemeComponent, MessageComponent, EmailComponent, RegistrySettingComponent],
    imports: [
        RouterModule,
        CoreModule,
        CredentialModule,
        BackupAccountModule,
    ]
})

export class SettingModule {

}

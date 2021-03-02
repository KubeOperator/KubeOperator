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
import { RegistryCreateComponent } from './registry-setting/registry-create/registry-create.component';
import { RegistryListComponent } from './registry-setting/registry-list/registry-list.component';
import { RegistryDeleteComponent } from './registry-setting/registry-delete/registry-delete.component';
import { RegistryUpdateComponent } from './registry-setting/registry-update/registry-update.component';
import {SharedModule} from '../../shared/shared.module';


@NgModule({
    declarations: [SettingComponent, SystemComponent, LicenseComponent, LicenseImportComponent,
        LdapComponent, ThemeComponent, MessageComponent, EmailComponent, RegistrySettingComponent, RegistryCreateComponent, RegistryListComponent, RegistryDeleteComponent, RegistryUpdateComponent],
    imports: [
        RouterModule,
        CoreModule,
        CredentialModule,
        BackupAccountModule,
        SharedModule,
    ]
})

export class SettingModule {

}

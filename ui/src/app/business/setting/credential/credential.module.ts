import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {CredentialDeleteComponent} from './credential-delete/credential-delete.component';
import {CoreModule} from '../../../core/core.module';
import {CredentialEditComponent} from './credential-edit/credential-edit.component';
import {CredentialListComponent} from './credential-list/credential-list.component';
import {CredentialComponent} from './credential.component';
import {CredentialCreateComponent} from './credential-create/credential-create.component';


@NgModule({
    declarations: [CredentialDeleteComponent, CredentialEditComponent, CredentialListComponent, CredentialComponent, CredentialCreateComponent],
    imports: [
        CommonModule,
        CoreModule
    ]
})
export class CredentialModule {
}

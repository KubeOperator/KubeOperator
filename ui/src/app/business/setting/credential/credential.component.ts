import {Component, OnInit, ViewChild} from '@angular/core';
import {CredentialListComponent} from './credential-list/credential-list.component';
import {CredentialCreateComponent} from './credential-create/credential-create.component';
import {CredentialDeleteComponent} from './credential-delete/credential-delete.component';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {Credential} from './credential';
import {CredentialService} from './credential.service';
import {CredentialEditComponent} from './credential-edit/credential-edit.component';

@Component({
    selector: 'app-credential',
    templateUrl: './credential.component.html',
    styleUrls: ['./credential.component.css']
})
export class CredentialComponent extends BaseModelDirective<Credential> implements OnInit {


    @ViewChild(CredentialListComponent, {static: true})
    list: CredentialListComponent;

    @ViewChild(CredentialCreateComponent, {static: true})
    create: CredentialCreateComponent;

    @ViewChild(CredentialDeleteComponent, {static: true})
    delete: CredentialDeleteComponent;

    @ViewChild(CredentialEditComponent, {static: true})
    edit: CredentialEditComponent;

    constructor(private credentialService: CredentialService) {
        super(credentialService);
    }

    ngOnInit(): void {
    }

    openCreate() {
        this.create.open();
    }

    openDelete(items: Credential[]) {
        this.delete.open(items);
    }

    openEdit(item: Credential) {
        this.edit.open(item);
    }

    refresh() {
        this.list.reset();
        this.list.refresh();
    }

}

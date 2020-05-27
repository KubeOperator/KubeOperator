import {Component, OnInit, ViewChild} from '@angular/core';
import {CredentialListComponent} from './credential-list/credential-list.component';
import {CredentialCreateComponent} from './credential-create/credential-create.component';
import {CredentialDeleteComponent} from './credential-delete/credential-delete.component';
import {BaseModelComponent} from '../../../shared/class/BaseModelComponent';
import {Credential} from './credential';
import {CredentialService} from './credential.service';

@Component({
    selector: 'app-credential',
    templateUrl: './credential.component.html',
    styleUrls: ['./credential.component.css']
})
export class CredentialComponent extends BaseModelComponent<Credential> implements OnInit {


    @ViewChild(CredentialListComponent, {static: true})
    list: CredentialListComponent;

    @ViewChild(CredentialCreateComponent, {static: true})
    create: CredentialCreateComponent;

    @ViewChild(CredentialDeleteComponent, {static: true})
    delete: CredentialDeleteComponent;


    constructor(private credentialService: CredentialService) {
        super(credentialService);
    }

    ngOnInit(): void {
    }

    openCreate() {
        this.create.open();
    }

    openDelete() {

    }

    refresh() {
        this.list.reset();
        this.list.refresh();
    }
}

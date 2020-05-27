import {Component, OnInit} from '@angular/core';
import {CredentialService} from '../credential.service';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {Credential} from '../credential';

@Component({
    selector: 'app-credential-list',
    templateUrl: './credential-list.component.html',
    styleUrls: ['./credential-list.component.css']
})
export class CredentialListComponent extends BaseModelComponent<Credential> implements OnInit {

    constructor(private credentialService: CredentialService) {
        super(credentialService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }

    onCreate() {

    }
}

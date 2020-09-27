import {Component, OnInit} from '@angular/core';
import {CredentialService} from '../credential.service';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {Credential} from '../credential';

@Component({
    selector: 'app-credential-list',
    templateUrl: './credential-list.component.html',
    styleUrls: ['./credential-list.component.css']
})
export class CredentialListComponent extends BaseModelDirective<Credential> implements OnInit {

    constructor(private credentialService: CredentialService) {
        super(credentialService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }
}

import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {CredentialService} from '../credential.service';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {Credential} from '../credential';
import {NgForm} from '@angular/forms';

@Component({
    selector: 'app-credential-edit',
    templateUrl: './credential-edit.component.html',
    styleUrls: ['./credential-edit.component.css']
})
export class CredentialEditComponent implements OnInit {

    item: Credential = new Credential();
    opened = false;
    isSubmitGoing = false;
    @ViewChild('credentialEditForm') credentialForm: NgForm;
    @Output() edit = new EventEmitter();

    constructor(private service: CredentialService) {
    }

    ngOnInit(): void {
    }

    open(item: Credential) {
        this.item = item;
        this.opened = true;
    }

    onCancel() {
        this.item = new Credential();
        this.opened = false;
    }

    onSubmit() {
        this.isSubmitGoing = true;
        this.service.update(this.item.name, this.item).subscribe(data => {
            this.opened = false;
            this.isSubmitGoing = false;
            this.edit.emit();
        });
    }
}

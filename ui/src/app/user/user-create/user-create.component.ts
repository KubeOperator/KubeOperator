import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {FormControl, FormGroup, NgForm} from '@angular/forms';
import {Host} from '../../host/host';
import {Credential} from '../../credential/credential-list/credential';
import * as globals from '../../globals';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {User} from '../user';

@Component({
    selector: 'app-user-create',
    templateUrl: './user-create.component.html',
    styleUrls: ['./user-create.component.css']
})
export class UserCreateComponent implements OnInit {
    @Output() create = new EventEmitter<boolean>();
    staticBackdrop = true;
    closable = false;
    opened: boolean;
    isSubmitGoing = false;
    user: User = new User();
    loading = false;
    @ViewChild('userForm', {static: true}) hostFrom: NgForm;

    ngOnInit() {

    }

    reset() {
    }


    onCancel() {
        this.opened = false;
    }

    onSubmit() {
    }

    newUser() {
        this.opened = true;
        this.reset();
        this.user = new User();
    }
}

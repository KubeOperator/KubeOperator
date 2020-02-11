import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {FormControl, FormGroup, NgForm} from '@angular/forms';
import {Host} from '../../host/host';
import {Credential} from '../../credential/credential-list/credential';
import * as globals from '../../globals';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {User} from '../user';
import {UserService} from '../user.service';

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
    @ViewChild('userForm', {static: true}) userFrom: NgForm;
    isPasswordMatch = true;
    isUserNameDuplicate = false;

    constructor(private userService: UserService) {
    }


    ngOnInit() {

    }

    reset() {
        this.isPasswordMatch = true;
        this.isUserNameDuplicate = false;
        this.userFrom.resetForm();
    }


    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        if (this.isSubmitGoing) {
            return;
        }
        this.isSubmitGoing = true;
        this.userService.createUser(this.user).subscribe(data => {
            this.isSubmitGoing = false;
            this.opened = false;
            this.create.emit(true);
        });
    }

    checkPassword() {
        this.isPasswordMatch = this.user.password === this.user.ensurePassword;
    }

    checkUsernameDuplicate() {
        this.userService.listUsers().subscribe(data => {
            data.some(u => {
                if (u.username === this.user.username) {
                    this.isUserNameDuplicate = true;
                    return;
                }
            });
        });
    }

    newUser() {
        this.reset();
        this.opened = true;
        this.user = new User();
    }
}

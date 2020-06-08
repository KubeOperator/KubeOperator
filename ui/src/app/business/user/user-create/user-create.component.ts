import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {User, UserCreateRequest} from '../user';
import {BaseModelComponent} from '../../../shared/class/BaseModelComponent';
import {UserService} from '../user.service';
import {passwordPattern} from '../../../globals';
import {NgForm} from '@angular/forms';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../shared/common-component/modal-alert/modal-alert.service';

@Component({
    selector: 'app-user-create',
    templateUrl: './user-create.component.html',
    styleUrls: ['./user-create.component.css']
})
export class UserCreateComponent extends BaseModelComponent<User> implements OnInit {

    opened = false;
    isSubmitGoing = false;
    item: UserCreateRequest = new UserCreateRequest();
    passwordPattern = passwordPattern;

    @ViewChild('userForm') userForm: NgForm;
    @Output() created = new EventEmitter();


    constructor(private userService: UserService, private modalAlertService: ModalAlertService) {
        super(userService);
    }

    ngOnInit(): void {
    }

    open() {
        this.opened = true;
        this.item = new UserCreateRequest();
        this.userForm.resetForm();
    }


    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.isSubmitGoing = true;
        this.userService.create(this.item).subscribe(data => {
            this.opened = false;
            this.isSubmitGoing = false;
            this.created.emit();
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}

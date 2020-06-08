import {Component, OnInit, ViewChild} from '@angular/core';
import {UserListComponent} from './user-list/user-list.component';
import {UserCreateComponent} from './user-create/user-create.component';
import {UserUpdateComponent} from './user-update/user-update.component';
import {UserDeleteComponent} from './user-delete/user-delete.component';

@Component({
    selector: 'app-user',
    templateUrl: './user.component.html',
    styleUrls: ['./user.component.css']
})
export class UserComponent implements OnInit {

    @ViewChild(UserListComponent, {static: true})
    list: UserListComponent;

    @ViewChild(UserCreateComponent, {static: true})
    create: UserCreateComponent;

    @ViewChild(UserUpdateComponent, {static: true})
    update: UserUpdateComponent;

    @ViewChild(UserDeleteComponent, {static: true})
    delete: UserDeleteComponent;

    constructor() {
    }

    ngOnInit(): void {
    }


    refresh() {
        this.list.reset();
        this.list.refresh();
    }

    openCreate() {
        this.create.open();
    }

    openDelete(items) {
        this.delete.open(items);
    }

    openUpdate(item) {
        this.update.open(item);
    }

}

import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {User} from '../user';
import {UserService} from '../user.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {CommonAlertService} from '../../base/header/common-alert.service';

@Component({
    selector: 'app-user-list',
    templateUrl: './user-list.component.html',
    styleUrls: ['./user-list.component.css']
})
export class UserListComponent implements OnInit {

    loading = true;
    users: User[] = [];
    selected: User[] = [];
    deleteModal = false;
    @Output() addUser = new EventEmitter();

    constructor(private userService: UserService, private alertService: CommonAlertService) {
    }

    ngOnInit() {
        this.listUser();
    }


    toggleActiveUser(user: User) {
        this.userService.activeUser(user).subscribe(data => {
            user = data;
        });
    }

    listUser() {
        this.loading = true;
        this.userService.listUsers().subscribe(data => {
            this.users = data.filter((u) => {
                return u.username !== 'admin';
            });
            this.loading = false;
        });
    }

    onDelete() {
        this.deleteModal = true;
    }

    confirmDelete() {
        const promises: Promise<{}>[] = [];
        this.selected.forEach(user => {
            promises.push(this.userService.deleteUser(user.id).toPromise());
        });
        Promise.all(promises).then(() => {
            this.refresh();
            this.alertService.showAlert('删除用户成功！', AlertLevels.SUCCESS);
        }).finally(() => {
            this.deleteModal = false;
            this.selected = [];
        });
    }

    addNewUser() {
        this.addUser.emit();
    }

    refresh() {
        this.listUser();
    }
}

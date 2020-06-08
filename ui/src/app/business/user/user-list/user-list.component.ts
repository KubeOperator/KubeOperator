import {Component, OnInit} from '@angular/core';
import {BaseModelComponent} from '../../../shared/class/BaseModelComponent';
import {User} from '../user';
import {UserService} from '../user.service';

@Component({
    selector: 'app-user-list',
    templateUrl: './user-list.component.html',
    styleUrls: ['./user-list.component.css']
})
export class UserListComponent extends BaseModelComponent<User> implements OnInit {

    constructor(private userService: UserService) {
        super(userService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }

}

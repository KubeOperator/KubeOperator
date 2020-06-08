import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelComponent} from '../../../shared/class/BaseModelComponent';
import {User} from '../user';
import {UserService} from '../user.service';

@Component({
    selector: 'app-user-delete',
    templateUrl: './user-delete.component.html',
    styleUrls: ['./user-delete.component.css']
})
export class UserDeleteComponent extends BaseModelComponent<User> implements OnInit {

    opened = false;
    items: User[] = [];

    @Output()
    deleted = new EventEmitter();

    constructor(private userService: UserService) {
        super(userService);
    }

    ngOnInit(): void {
    }

    open(items) {
        this.items = items;
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.service.batch('delete', this.items).subscribe(data => {
            this.deleted.emit();
            this.opened = false;
        });
    }
}

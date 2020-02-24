import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {User} from '../../../../user/user';
import {UserService} from '../../../../user/user.service';
import {SessionService} from '../../../../shared/session.service';
import {NgForm} from '@angular/forms';

@Component({
  selector: 'app-item-change',
  templateUrl: './item-change.component.html',
  styleUrls: ['./item-change.component.css']
})
export class ItemChangeComponent implements OnInit {
  user: User;
  isSubmitGoing = false;
  opened = false;
  loading = true;
  @ViewChild('itemForm', {static: true})
  itemForm: NgForm;
  @Output() changeItem = new EventEmitter();

  constructor(private userService: UserService, private sessionService: SessionService) {
  }

  ngOnInit() {
  }

  open() {
    this.itemForm.resetForm();
    this.opened = true;
    this.loading = true;
    const sessionUser = this.sessionService.getCacheUser();
    this.userService.getUser(sessionUser.id).subscribe(data => {
      this.user = data;
      this.loading = false;
    });
  }

  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.isSubmitGoing = true;
    this.userService.updateUser(this.user).subscribe(() => {
      this.changeItem.emit();
      this.opened = false;
      this.isSubmitGoing = false;
    });
  }

  onCancel() {
    this.opened = false;
  }

}

import {AfterViewChecked, Component, Input, OnInit, ViewChild} from '@angular/core';
import {NgForm} from '@angular/forms';
import {SignInCredential} from '../../shared/signInCredential';
import {ActivatedRoute, Router} from '@angular/router';
import {SessionService} from '../../shared/session.service';
import {CommonRoutes} from '../../shared/shared.const';
import {MessageService} from '../../base/message.service';
import {MessageLevels} from '../../base/message/message-level';
import {SettingService} from '../../setting/setting.service';


export const signInStatusNormal = 0;
export const signInStatusOnGoing = 1;
export const signInStatusError = -1;

@Component({
  selector: 'app-sign-in',
  templateUrl: './sign-in.component.html',
  styleUrls: ['./sign-in.component.css']
})
export class SignInComponent implements OnInit, AfterViewChecked {
  redireUrl = '';
  message = '';

  // form
  signInFrom: NgForm;
  @ViewChild('signInForm', { static: true }) currentForm: NgForm;

  signInStatus: number = signInStatusNormal;

  @Input() signInCredential: SignInCredential = {
    principal: '',
    password: ''
  };

  constructor(private router: Router, private  route: ActivatedRoute, private session: SessionService, private settingService: SettingService, private messageService: MessageService) {
  }

  ngAfterViewChecked(): void {
    if (this.signInStatus === signInStatusError) {
      this.formChanged();
    }
  }

  ngOnInit(): void {
  }

  public get isError(): boolean {
    return this.signInStatus === signInStatusError;
  }

  public get isOnGoing(): boolean {
    return this.signInStatus === signInStatusOnGoing;
  }

  public get isValid(): boolean {
    return this.currentForm.form.valid;
  }


  handleError(error: any) {
    this.signInStatus = signInStatusError;
    if (error.status === 504 || error.status === 502) {
      this.message = 'kubeOperator Api 连接失败！';
    } else if (error.status === 400) {
      this.message = '用户名或密码错误！';
    }
  }

  formChanged() {
    if (this.currentForm === this.signInFrom) {
      return;
    }
    this.signInFrom = this.currentForm;
    if (this.signInFrom) {
      this.signInFrom.valueChanges.subscribe((data) => {
        this.updateState();
      });
    }
  }


  updateState(): void {
    if (this.signInStatus === signInStatusError) {
      this.signInStatus = signInStatusNormal; // reset
    }
  }


  signIn(): void {
    if (!this.isValid) {
      this.signInStatus = signInStatusError;
      return;
    }

    if (this.isOnGoing) {
      return;
    }

    this.signInStatus = signInStatusOnGoing;
    this.session.authUser(this.signInCredential).subscribe(data => {
      this.session.cacheToken(data);
      if (this.redireUrl === '') {
        this.router.navigateByUrl(CommonRoutes.F2O_DEFAULT);
      } else {
        this.router.navigateByUrl(this.redireUrl);
      }
    }, (error) => this.handleError(error));

  }

  // 设置主机名


}

import {AfterViewChecked, Component, Input, OnInit, ViewChild} from '@angular/core';
import {NgForm} from '@angular/forms';
import {SignInCredential} from '../../shared/signInCredential';
import {ActivatedRoute, Router} from '@angular/router';
import {SessionService} from '../../shared/session.service';
import {CommonRoutes} from '../../shared/shared.const';
import {SessionUser} from '../../shared/session-user';


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


  // form
  signInFrom: NgForm;
  @ViewChild('signInForm') currentForm: NgForm;

  signInStatus: number = signInStatusNormal;

  @Input() signInCredential: SignInCredential = {
    principal: '',
    password: ''
  };

  constructor(private router: Router, private  route: ActivatedRoute, private session: SessionService) {
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
    const message = error.status ? error.status + ':' + error.statusText : error;
    console.error('An error occurred when signing in:', message);
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
      this.cacheToken(data);
      if (this.redireUrl === '') {
        this.router.navigateByUrl(CommonRoutes.F2O_DEFAULT);
      } else {
        this.router.navigateByUrl(this.redireUrl);
      }
    }, (error) => this.handleError(error));

  }

  cacheToken(user: SessionUser) {
    localStorage.setItem('current_user', JSON.stringify(user));
  }

}

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CredentialEditComponent } from './credential-edit.component';

describe('CredentialEditComponent', () => {
  let component: CredentialEditComponent;
  let fixture: ComponentFixture<CredentialEditComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CredentialEditComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CredentialEditComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

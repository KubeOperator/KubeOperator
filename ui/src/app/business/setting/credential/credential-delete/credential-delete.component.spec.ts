import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CredentialDeleteComponent } from './credential-delete.component';

describe('CredentialDeleteComponent', () => {
  let component: CredentialDeleteComponent;
  let fixture: ComponentFixture<CredentialDeleteComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CredentialDeleteComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CredentialDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

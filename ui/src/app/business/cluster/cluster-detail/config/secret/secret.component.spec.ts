import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SecretComponent } from './secret.component';

describe('SecretComponent', () => {
  let component: SecretComponent;
  let fixture: ComponentFixture<SecretComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SecretComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SecretComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

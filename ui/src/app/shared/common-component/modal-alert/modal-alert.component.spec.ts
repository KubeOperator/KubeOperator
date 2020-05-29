import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ModalAlertComponent } from './modal-alert.component';

describe('ModalAlertComponent', () => {
  let component: ModalAlertComponent;
  let fixture: ComponentFixture<ModalAlertComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ModalAlertComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ModalAlertComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

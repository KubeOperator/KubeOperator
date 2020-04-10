import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OpsCreateComponent } from './ops-create.component';

describe('OpsCreateComponent', () => {
  let component: OpsCreateComponent;
  let fixture: ComponentFixture<OpsCreateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ OpsCreateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OpsCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

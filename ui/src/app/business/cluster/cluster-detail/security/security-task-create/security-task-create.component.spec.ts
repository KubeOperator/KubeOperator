import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SecurityTaskCreateComponent } from './security-task-create.component';

describe('SecurityTaskCreateComponent', () => {
  let component: SecurityTaskCreateComponent;
  let fixture: ComponentFixture<SecurityTaskCreateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SecurityTaskCreateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SecurityTaskCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

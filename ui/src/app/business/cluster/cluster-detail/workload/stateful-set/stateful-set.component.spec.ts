import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StatefulSetComponent } from './stateful-set.component';

describe('StatefulSetComponent', () => {
  let component: StatefulSetComponent;
  let fixture: ComponentFixture<StatefulSetComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StatefulSetComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StatefulSetComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

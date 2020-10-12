import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VmConfigComponent } from './vm-config.component';

describe('VmConfigComponent', () => {
  let component: VmConfigComponent;
  let fixture: ComponentFixture<VmConfigComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ VmConfigComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(VmConfigComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
